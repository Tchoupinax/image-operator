package controller

import (
	"context"
	"fmt"
	"os"
	"strings"

	skopeoiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/v1alpha1"
	"github.com/Tchoupinax/skopeo-operator/internal/helpers"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateSkopeoJob(
	r *ImageReconciler,
	ctx context.Context,
	req ctrl.Request,
	image skopeoiov1alpha1.Image,
	logger logr.Logger,
	overridenVersion string,
) {
	desiredJob := GenerateSkopeoJob(
		r,
		ctx,
		req,
		image,
		logger,
		overridenVersion,
	)

	creationError := r.Create(ctx, &desiredJob)
	if creationError != nil {
		fmt.Println(creationError)
	} else {
		// If the creation is a success, increase metrics
		var imageKey = fmt.Sprintf("%s:%s", image.Spec.Source.ImageName, image.Spec.Source.ImageVersion)
		r.PrometheusReloadGauge.WithLabelValues(imageKey).Inc()
		r.LastTimeImageWasReloaded.WithLabelValues(imageKey).Set(helpers.GenerateTimestamp())
	}
}

func GenerateSkopeoJob(
	r *ImageReconciler,
	ctx context.Context,
	req ctrl.Request,
	image skopeoiov1alpha1.Image,
	logger logr.Logger,
	overridenVersion string,
) batchv1.Job {
	logger.Info("Create job to copy image", image.Spec.Source.ImageName, image.Spec.Source.ImageVersion)

	skopeoCommand := []string{
		"skopeo",
		"copy",
		fmt.Sprintf("docker://%s:%s", image.Spec.Source.ImageName, overridenVersion),
		fmt.Sprintf("docker://%s:%s", image.Spec.Destination.ImageName, overridenVersion),
		"--all",
		"--preserve-digests",
	}

	if os.Getenv("CREDS_DESTINATION_USERNAME") != "" && os.Getenv("CREDS_DESTINATION_PASSWORD") != "" {
		skopeoCommand = append(
			skopeoCommand,
			fmt.Sprintf("--dest-creds=%s:%s", os.Getenv("CREDS_DESTINATION_USERNAME"), os.Getenv("CREDS_DESTINATION_PASSWORD")),
		)
	}

	if os.Getenv("CREDS_SOURCE_USERNAME") != "" && os.Getenv("CREDS_SOURCE_PASSWORD") != "" {
		skopeoCommand = append(
			skopeoCommand,
			fmt.Sprintf("--dest-creds=%s:%s", os.Getenv("CREDS_SOURCE_USERNAME"), os.Getenv("CREDS_SOURCE_PASSWORD")),
		)
	}

	if os.Getenv("DISABLE_SRC_TLS_VERIFICATION") == "true" {
		skopeoCommand = append(
			skopeoCommand,
			"--src-tls-verify=false",
		)
	} else {
		skopeoCommand = append(
			skopeoCommand,
			"--src-tls-verify=true",
		)
	}

	if os.Getenv("DISABLE_DEST_TLS_VERIFICATION") == "true" {
		skopeoCommand = append(
			skopeoCommand,
			"--dest-tls-verify=false",
		)
	} else {
		skopeoCommand = append(
			skopeoCommand,
			"--dest-tls-verify=true",
		)
	}

	podNamespace := os.Getenv("PULL_JOB_NAMESPACE")
	if podNamespace == "" {
		podNamespace = "skopeo-operator"
	}

	skopeoImage := os.Getenv("SKOPEO_IMAGE")
	if skopeoImage == "" {
		skopeoImage = "quay.io/containers/skopeo"
	}

	skopeoVersion := os.Getenv("SKOPEO_VERSION")
	if skopeoVersion == "" {
		skopeoVersion = "v1.16.1"
	}

	var arguments []string
	if image.Spec.Destination.UseAwsIRSA || image.Spec.Source.UseAwsIRSA {
		arguments = []string{
			"-c",
			fmt.Sprintf(`yum install -y awscli &&
			aws ecr get-login-password --region eu-west-1 | skopeo login --username AWS --password-stdin 326954429656.dkr.ecr.eu-west-1.amazonaws.com &&
			%s`, strings.Join(skopeoCommand, " ")),
		}
	} else {
		arguments = []string{
			"-c",
			strings.Join(skopeoCommand, " "),
		}
	}

	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf(
				"skopeo-job-copy-%s-%s",
				strings.ReplaceAll(req.Name, ".", ""),
				strings.ReplaceAll(overridenVersion, ".", ""),
			),
			Namespace: podNamespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            int32Ptr(0),
			TTLSecondsAfterFinished: int32Ptr(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: image.Annotations,
					Labels:      image.Labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					Containers: []corev1.Container{
						{
							Image:   fmt.Sprintf("%s:%s", skopeoImage, skopeoVersion),
							Name:    "skopeo",
							Command: []string{"/bin/bash"},
							Args:    arguments,
						},
					},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
