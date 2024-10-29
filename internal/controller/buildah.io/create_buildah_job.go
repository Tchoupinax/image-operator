package buildahio

import (
	"context"
	"fmt"
	"os"
	"strings"

	buildahiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/buildah.io/v1alpha1"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateBuildahJobs(
	r *ImageBuilderReconciler,
	ctx context.Context,
	req ctrl.Request,
	imageBuilder buildahiov1alpha1.ImageBuilder,
	logger logr.Logger,
) (ctrl.Result, error) {
	jobNamespace := os.Getenv("BUILD_JOB_NAMESPACE")
	if jobNamespace == "" {
		jobNamespace = "skopeo-operator"
	}

	architectures := []buildahiov1alpha1.Architecture{imageBuilder.Spec.Architecture}
	if imageBuilder.Spec.Architecture == buildahiov1alpha1.Both {
		architectures = []buildahiov1alpha1.Architecture{buildahiov1alpha1.AMD64, buildahiov1alpha1.ARM64}
	}

	for _, arch := range architectures {
		jobName := fmt.Sprintf(
			"buildah-job-%s-%s",
			req.Name,
			string(arch),
		)

		configMap := generateConfigMap(imageBuilder, jobName, jobNamespace)
		// https://github.com/kubernetes-sigs/kubebuilder/issues/892#issuecomment-698022177
		if referenceError := controllerutil.SetControllerReference(&imageBuilder, &configMap, r.Scheme); referenceError != nil {
			logger.Error(referenceError, "Fail to set the reference on the configmap")
			return ctrl.Result{}, referenceError
		}

		if createError := r.Create(ctx, &configMap); createError != nil {
			logger.Error(createError, "Fail to create config map")
			return ctrl.Result{}, createError
		}

		desiredJob := GenerateBuildahJob(
			r,
			ctx,
			req,
			imageBuilder,
			logger,
			jobName,
			jobNamespace,
			arch,
		)
		// https://github.com/kubernetes-sigs/kubebuilder/issues/892#issuecomment-698022177
		if referenceError2 := controllerutil.SetControllerReference(&imageBuilder, &desiredJob, r.Scheme); referenceError2 != nil {
			logger.Error(referenceError2, "Fail to set the reference on the job")
			return ctrl.Result{}, referenceError2
		}

		if createError2 := r.Create(ctx, &desiredJob); createError2 != nil {
			logger.Error(createError2, "Fail to create job")
			return ctrl.Result{}, createError2
		}

		var imageBuilderLocal buildahiov1alpha1.ImageBuilder
		if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &imageBuilderLocal); err != nil {
			logger.Info(fmt.Sprintf("Object %s not found (%s)", req.Name, req.Namespace))
		}

		imageBuilderLocal.Status.RanJobs = append(imageBuilderLocal.Status.RanJobs, jobName)

		updateError := r.Status().Update(ctx, &imageBuilderLocal)
		if updateError != nil {
			logger.Error(updateError, "Error when updating Image Builder resource")
		}

	}

	return ctrl.Result{}, nil
}

func generateConfigMap(
	imageBuilder buildahiov1alpha1.ImageBuilder,
	jobName string,
	jobNamespace string,
) corev1.ConfigMap {
	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},
		Data: map[string]string{
			"dockerfile": imageBuilder.Spec.Source,
		},
	}
}

func GenerateBuildahJob(
	r *ImageBuilderReconciler,
	ctx context.Context,
	req ctrl.Request,
	imageBuilder buildahiov1alpha1.ImageBuilder,
	logger logr.Logger,
	jobName string,
	jobNamespace string,
	architecture buildahiov1alpha1.Architecture,
) batchv1.Job {
	tag := fmt.Sprintf("%s-%s", imageBuilder.Spec.Image.ImageVersion, string(architecture))
	imageName := fmt.Sprintf("%s:%s", imageBuilder.Spec.Image.ImageName, tag)

	logger.Info(fmt.Sprintf("GenerateBuildahJob %s", imageName))

	var arguments []string
	if imageBuilder.Spec.Image.UseAwsIRSA {
		parts := strings.Split(imageBuilder.Spec.Image.ImageName, "/")
		arguments = []string{
			fmt.Sprintf("echo '{\"credHelpers\":{\"%s\":\"ecr-login\"}}' > /tmp/auth.json && buildah build -t %s --storage-driver=vfs -f /docker/Dockerfile && buildah push --storage-driver=vfs %s", parts[0], imageName, imageName),
		}
	} else {
		var username = os.Getenv("CREDS_DESTINATION_USERNAME")
		var password = os.Getenv("CREDS_DESTINATION_PASSWORD")

		arguments = []string{
			fmt.Sprintf("buildah build -t %s --storage-driver=vfs -f /docker/Dockerfile && buildah push --storage-driver=vfs --creds=%s:%s %s", imageName, username, password, imageName),
		}
	}

	job := GenerateAbtractBuildahJob(
		imageBuilder,
		string(architecture),
		arguments,
	)

	job.Spec.Template.Spec.Affinity = &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      "kubernetes.io/arch",
						Operator: "In",
						Values:   []string{string(architecture)},
					}},
				}},
			},
		},
	}

	job.Spec.Template.Spec.Volumes = []corev1.Volume{{
		Name: "config-volume",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: jobName,
				},
				Items: []corev1.KeyToPath{
					{
						Key:  "dockerfile",
						Path: "Dockerfile",
					},
				},
			},
		},
	}}

	job.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{{
		Name:      "config-volume",
		MountPath: "/docker",
	}}

	return job
}
