package controller

import (
	"context"
	"os"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	phaseRunning   = "RUNNING"
	phaseCompleted = "COMPLETED"
	phaseFailed    = "FAILED"
)

func pullJobNamespace() string {
	podNamespace := os.Getenv("PULL_JOB_NAMESPACE")
	if podNamespace == "" {
		return "image-operator"
	}
	return podNamespace
}

func isJobComplete(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobFailed(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func versionsToMonitor(
	logger logr.Logger,
	image *skopeoiov1alpha1.Image,
) []string {
	selectedVersions := listVersionsForImage(
		logger,
		image.Spec.Source.ImageName,
		image.Spec.Source.ImageVersion,
		image.Spec.AllowCandidateRelease,
		helpers.DockerHubAuth{},
		helpers.AWSPublicECR{},
	)

	if image.Spec.Mode == skopeoiov1alpha1.ONCE_BY_TAG {
		return helpers.Filter(selectedVersions, func(tag string) bool {
			return !helpers.Contains(image.Status.TagAlreadySynced, tag)
		})
	}

	return selectedVersions
}

func updateStatusFromJobs(
	r *ImageReconciler,
	ctx context.Context,
	req ctrl.Request,
	image *skopeoiov1alpha1.Image,
	logger logr.Logger,
) (bool, error) {
	if len(image.Status.History) == 0 {
		return false, nil
	}

	versions := versionsToMonitor(logger, image)
	if len(versions) == 0 && image.Spec.Mode == skopeoiov1alpha1.ONCE_BY_TAG {
		image.Status.Phase = phaseCompleted
		return true, r.Status().Update(ctx, image)
	}

	var (
		foundJobs   int
		runningJobs int
		failedJobs  int
		statusDirty bool
	)

	for _, version := range versions {
		var job batchv1.Job
		err := r.Get(ctx, types.NamespacedName{
			Namespace: pullJobNamespace(),
			Name:      helpers.GenerateSkopeoJobName(req.Name, version),
		}, &job)
		if apierrors.IsNotFound(err) {
			continue
		}
		if err != nil {
			return false, err
		}

		foundJobs++

		switch {
		case isJobFailed(&job):
			failedJobs++
		case isJobComplete(&job):
			if image.Spec.Mode == skopeoiov1alpha1.ONCE_BY_TAG &&
				!helpers.Contains(image.Status.TagAlreadySynced, version) {
				image.Status.TagAlreadySynced = append(image.Status.TagAlreadySynced, version)
				statusDirty = true
			}
		default:
			runningJobs++
		}
	}

	nextPhase := image.Status.Phase
	switch {
	case failedJobs > 0:
		nextPhase = phaseFailed
	case foundJobs > 0 && runningJobs == 0 && failedJobs == 0:
		nextPhase = phaseCompleted
	case foundJobs > 0:
		nextPhase = phaseRunning
	case len(versions) > 0 && image.Status.Phase == "":
		nextPhase = phaseRunning
	}

	if nextPhase != image.Status.Phase {
		image.Status.Phase = nextPhase
		statusDirty = true
	}

	if !statusDirty {
		return false, nil
	}

	return true, r.Status().Update(ctx, image)
}

func persistImageStatus(
	r *ImageReconciler,
	ctx context.Context,
	image *skopeoiov1alpha1.Image,
) error {
	return r.Status().Update(ctx, image)
}
