package jobgen_test

import (
	"testing"

	buildahiov1alpha1 "github.com/Tchoupinax/image-operator/api/buildah.io/v1alpha1"
	buildahio "github.com/Tchoupinax/image-operator/internal/controller/buildah.io"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestGenerateBuildahManifestJobUsesCredentials(t *testing.T) {
	t.Setenv("CREDS_DESTINATION_USERNAME", "dest-user")
	t.Setenv("CREDS_DESTINATION_PASSWORD", "dest-pass")

	job := buildahio.GenerateBuildahManifestJob(
		&buildahio.ImageBuilderReconciler{},
		nil,
		ctrl.Request{},
		buildahiov1alpha1.ImageBuilder{
			Spec: buildahiov1alpha1.ImageBuilderSpec{
				Image: buildahiov1alpha1.ImageEndpoint{
					ImageName:    "registry.example.com/app",
					ImageVersion: "1.2.3",
				},
			},
		},
		logr.Discard(),
	)

	if len(job.Spec.Template.Spec.Containers) != 1 {
		t.Fatalf("expected one container, got %d", len(job.Spec.Template.Spec.Containers))
	}

	expected := "buildah manifest create registry.example.com/app:1.2.3 && buildah manifest add --creds=dest-user:dest-pass --arch arm64 registry.example.com/app:1.2.3 registry.example.com/app:1.2.3-arm64 && buildah manifest add --creds=dest-user:dest-pass --arch amd64 registry.example.com/app:1.2.3 registry.example.com/app:1.2.3-amd64 && buildah manifest --creds=dest-user:dest-pass push registry.example.com/app:1.2.3"
	if job.Spec.Template.Spec.Containers[0].Args[0] != expected {
		t.Fatalf("unexpected manifest command:\n%s", job.Spec.Template.Spec.Containers[0].Args[0])
	}
}

func TestGenerateAbtractBuildahJobUsesNamespace(t *testing.T) {
	t.Setenv("BUILDAH_JOB_NAMESPACE", "buildah-jobs")

	job := buildahio.GenerateAbtractBuildahJob(
		buildahiov1alpha1.ImageBuilder{Name: "demo"},
		"amd64",
		[]string{"buildah bud ."},
	)

	if job.Namespace != "buildah-jobs" {
		t.Fatalf("expected namespace buildah-jobs, got %s", job.Namespace)
	}
	if job.Spec.Template.Spec.Containers[0].Args[0] != "buildah bud ." {
		t.Fatalf("unexpected job args: %v", job.Spec.Template.Spec.Containers[0].Args)
	}
}
