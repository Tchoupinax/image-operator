package buildahio

import (
	"context"
	"fmt"
	"os"
	"strings"

	buildahiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/buildah.io/v1alpha1"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateBuildahManifestJobs(
	r *ImageBuilderReconciler,
	ctx context.Context,
	req ctrl.Request,
	imageBuilder buildahiov1alpha1.ImageBuilder,
	logger logr.Logger,
) {
	imageBuilder.Status.Pushed = true

	updateError := r.Status().Update(ctx, &imageBuilder)
	if updateError != nil {
		fmt.Println(updateError)
	}

	desiredJob := GenerateBuildahManifestJob(
		r,
		ctx,
		req,
		imageBuilder,
		logger,
	)

	// https://github.com/kubernetes-sigs/kubebuilder/issues/892#issuecomment-698022177
	err := controllerutil.SetControllerReference(&imageBuilder, &desiredJob, r.Scheme)
	if err != nil {
		fmt.Println(err)
	}

	errCreate := r.Create(ctx, &desiredJob)
	if errCreate != nil {
		fmt.Println(errCreate)
	}
}

func GenerateBuildahManifestJob(
	r *ImageBuilderReconciler,
	ctx context.Context,
	req ctrl.Request,
	imageBuilder buildahiov1alpha1.ImageBuilder,
	logger logr.Logger,
) batchv1.Job {
	tag := imageBuilder.Spec.Image.ImageVersion
	imageNameArm := fmt.Sprintf("%s:%s-arm64", imageBuilder.Spec.Image.ImageName, tag)
	imageNameAmd := fmt.Sprintf("%s:%s-amd64", imageBuilder.Spec.Image.ImageName, tag)
	imageName := fmt.Sprintf("%s:%s", imageBuilder.Spec.Image.ImageName, tag)

	logger.Info(fmt.Sprintf("GenerateBuildahManifestJob %s", imageName))

	var arguments []string
	if imageBuilder.Spec.Image.UseAwsIRSA {
		parts := strings.Split(imageBuilder.Spec.Image.ImageName, "/")
		arguments = []string{
			fmt.Sprintf("echo '{\"credHelpers\":{\"%s\":\"ecr-login\"}}' > /tmp/auth.json && buildah manifest create %s && buildah manifest add --arch arm64 %s %s && buildah manifest add --arch amd64 %s %s && buildah manifest push %s", parts[0], imageName, imageName, imageNameArm, imageName, imageNameAmd, imageName),
		}
	} else {
		var username = os.Getenv("CREDS_DESTINATION_USERNAME")
		var password = os.Getenv("CREDS_DESTINATION_PASSWORD")

		create := fmt.Sprintf("buildah manifest create %s", imageName)
		addArm64 := fmt.Sprintf("buildah manifest add --creds=%s:%s --arch amd64 %s %s", username, password, imageName, imageNameAmd)
		addAmd64 := fmt.Sprintf("buildah manifest add --creds=%s:%s --arch arm64 %s %s", username, password, imageName, imageNameArm)
		push := fmt.Sprintf("buildah manifest --creds=%s:%s push %s", username, password, imageName)

		arguments = []string{
			fmt.Sprintf("%s && %s && %s && %s", create, addAmd64, addArm64, push),
		}
	}

	return GenerateAbtractBuildahJob(
		imageBuilder,
		"manifest",
		arguments,
	)
}
