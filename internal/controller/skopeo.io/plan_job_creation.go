package controller

import (
	"context"
	"time"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func planJobCreation(
	r *ImageReconciler,
	ctx context.Context,
	req ctrl.Request,
	image *skopeoiov1alpha1.Image,
	logger logr.Logger,
) ctrl.Result {
	addHistory(image)

	selectedVersions := helpers.ListVersion(
		logger,
		image.Spec.Source.ImageName,
		image.Spec.Source.ImageVersion,
		image.Spec.AllowCandidateRelease,
		helpers.DockerHubAuth{},
		helpers.AWSPublicECR{},
	)

	if image.Spec.Mode == skopeoiov1alpha1.ONCE_BY_TAG {
		selectedVersions = helpers.Filter(selectedVersions, func(tag string) bool {
			return !helpers.Contains(image.Status.TagAlreadySynced, tag)
		})
	}

	if len(selectedVersions) > 0 {
		logger.Info("Reload image")

		for _, selectedVersion := range selectedVersions {
			CreateSkopeoJob(r, ctx, req, image, logger, selectedVersion)
		}

		if image.Spec.Mode == "OnceByTag" {
			// This operation should only be done if the job succeeded
			image.Status.TagAlreadySynced = append(
				image.Status.TagAlreadySynced,
				selectedVersions...,
			)

			updateError := r.Status().Update(ctx, image)
			if updateError != nil {
				logger.Error(updateError, "Failed to update the CRD "+image.Name)
				// We stop the loop
				return ctrl.Result{}
			}
		}
	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 10 * time.Second,
	}
}

func addHistory(
	image *skopeoiov1alpha1.Image,
) {
	image.Status.History = append(image.Status.History, skopeoiov1alpha1.History{
		PerformedAt: metav1.Now(),
	})

	maxItems := 20
	if len(image.Status.History) > maxItems {
		image.Status.History = image.Status.History[len(image.Status.History)-maxItems:]
	}
}
