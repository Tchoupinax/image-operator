package controller

import (
	"context"
	"fmt"

	skopeoiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/v1alpha1"
	"github.com/Tchoupinax/skopeo-operator/internal/helpers"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func planJobCreation(
	r *ImageReconciler,
	ctx context.Context,
	req ctrl.Request,
	image skopeoiov1alpha1.Image,
	logger logr.Logger,
) {
	addHistory(r, ctx, image)

	selectedVersions := helpers.ListVersion(
		logger,
		image.Spec.Source.ImageName,
		image.Spec.Source.ImageVersion,
		image.Spec.AllowCandidateRelease,
		helpers.DockerHubAuth{},
		helpers.AWSPublicECR{},
	)

	if image.Spec.Mode == "OneShot" {
		selectedVersions = helpers.Filter(selectedVersions, func(tag string) bool {
			return !helpers.Contains(image.Status.TagAlreadySynced, tag)
		})
	}

	if len(selectedVersions) > 0 {
		logger.Info("Reload image")

		for _, tag := range selectedVersions {
			CreateSkopeoJob(r, ctx, req, image, logger, tag)
		}

		if image.Spec.Mode == "OneShot" {
			// This operation should only be done if the job succeeded
			image.Status.TagAlreadySynced = append(
				image.Status.TagAlreadySynced,
				selectedVersions...,
			)
			updateError2 := r.Status().Update(ctx, &image)
			if updateError2 != nil {
				fmt.Println(updateError2)
			}
		}
	}
}

func addHistory(
	r *ImageReconciler,
	ctx context.Context,
	image skopeoiov1alpha1.Image,
) {
	image.Status.History = append(image.Status.History, skopeoiov1alpha1.History{
		PerformedAt: metav1.Now(),
	})

	maxItems := 20
	if len(image.Status.History) > maxItems {
		image.Status.History = image.Status.History[len(image.Status.History)-maxItems:]
	}

	updateError := r.Status().Update(ctx, &image)
	if updateError != nil {
		fmt.Println(updateError)
	}
}
