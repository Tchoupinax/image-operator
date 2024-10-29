package buildahio

import (
	"fmt"
	"os"

	buildahiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/buildah.io/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateAbtractBuildahJob(
	imageBuilder buildahiov1alpha1.ImageBuilder,
	name string,
	arguments []string,
) batchv1.Job {
	jobNamespace := os.Getenv("BUILDAH_JOB_NAMESPACE")
	if jobNamespace == "" {
		jobNamespace = "skopeo-operator"
	}

	buildahImage := os.Getenv("BUILDAH_IMAGE")
	if buildahImage == "" {
		buildahImage = "quay.io/containers/buildah"
	}

	buildahVersion := os.Getenv("BUILDAH_VERSION")
	if buildahVersion == "" {
		buildahVersion = "v1.37.3"
	}

	privilegedContainer := false
	if os.Getenv("BUILDAH_PRIVILEGED_CONTAINER") != "" {
		if os.Getenv("BUILDAH_PRIVILEGED_CONTAINER") == "true" {
			privilegedContainer = true
		} else {
			privilegedContainer = false
		}
	}

	memoryLimit := "8Gi"
	if imageBuilder.Spec.Resources.Limits.Memory != "" {
		memoryLimit = imageBuilder.Spec.Resources.Limits.Memory
	}
	requestCpu := "200m"
	if imageBuilder.Spec.Resources.Requests.Cpu != "" {
		requestCpu = imageBuilder.Spec.Resources.Requests.Cpu
	}
	requestMemory := "200m"
	if imageBuilder.Spec.Resources.Requests.Memory != "" {
		requestMemory = imageBuilder.Spec.Resources.Requests.Memory
	}

	jobName := fmt.Sprintf("buildah-job-%s-%s", imageBuilder.Name, name)

	envs := []corev1.EnvVar{
		{
			Name:  "BUILDAH_ISOLATION",
			Value: "chroot",
		},
		{
			Name:  "BUILDAH_FORMAT",
			Value: "docker",
		},
		{
			Name:  "BUILDAH_LAYERS",
			Value: "true",
		},
	}

	if imageBuilder.Spec.Image.UseAwsIRSA {
		envs = append(envs, corev1.EnvVar{
			Name:  "REGISTRY_AUTH_FILE",
			Value: "/tmp/auth.json",
		})
	}

	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            int32Ptr(0),
			TTLSecondsAfterFinished: int32Ptr(3600), // We put 1 hour to have time to do the entire process before removing jobs
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: imageBuilder.Annotations,
					Labels:      imageBuilder.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "buildah",
							Image: fmt.Sprintf("%s:%s", buildahImage, buildahVersion),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"memory": resource.Quantity{
										Format: resource.Format(memoryLimit),
									},
								},
								Requests: corev1.ResourceList{
									"cpu": resource.Quantity{
										Format: resource.Format(requestCpu),
									},
									"memory": resource.Quantity{
										Format: resource.Format(requestMemory),
									},
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(privilegedContainer),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{
										"NET_RAW",
									},
								},
							},
							Command: []string{"/bin/sh", "-c"},
							Args:    arguments,
							Env:     envs,
						},
					},
					// PreemptionPolicy: corev1.PreemptionPolicy, // Find th correct import
					Priority:      int32Ptr(0),
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
