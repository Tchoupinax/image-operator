package controller

import (
	"context"
	"os"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status helpers", func() {
	Describe("pullJobNamespace", func() {
		It("defaults to image-operator when unset", func() {
			os.Unsetenv("PULL_JOB_NAMESPACE")
			Expect(pullJobNamespace()).To(Equal("image-operator"))
		})

		It("uses PULL_JOB_NAMESPACE when set", func() {
			Expect(os.Setenv("PULL_JOB_NAMESPACE", "custom-ns")).To(Succeed())
			DeferCleanup(func() {
				os.Unsetenv("PULL_JOB_NAMESPACE")
			})
			Expect(pullJobNamespace()).To(Equal("custom-ns"))
		})
	})

	Describe("isJobComplete", func() {
		It("returns true when the complete condition is true", func() {
			job := &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{{
						Type:   batchv1.JobComplete,
						Status: corev1.ConditionTrue,
					}},
				},
			}
			Expect(isJobComplete(job)).To(BeTrue())
		})

		It("returns false when the job is still running", func() {
			Expect(isJobComplete(&batchv1.Job{})).To(BeFalse())
		})
	})

	Describe("isJobFailed", func() {
		It("returns true when the failed condition is true", func() {
			job := &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{{
						Type:   batchv1.JobFailed,
						Status: corev1.ConditionTrue,
					}},
				},
			}
			Expect(isJobFailed(job)).To(BeTrue())
		})
	})
})

var _ = Describe("updateStatusFromJobs", func() {
	var (
		ctx              context.Context
		scheme           *runtime.Scheme
		cl               client.Client
		reconciler       *ImageReconciler
		originalLister   listVersionsFunc
		imageNamespaced  types.NamespacedName
		stubbedVersions  []string
	)

	BeforeEach(func() {
		ctx = context.Background()
		scheme = runtime.NewScheme()
		Expect(skopeoiov1alpha1.AddToScheme(scheme)).To(Succeed())
		Expect(batchv1.AddToScheme(scheme)).To(Succeed())

		cl = fake.NewClientBuilder().
			WithScheme(scheme).
			WithStatusSubresource(&skopeoiov1alpha1.Image{}).
			Build()

		reconciler = &ImageReconciler{
			Client: cl,
			Scheme: scheme,
		}

		imageNamespaced = types.NamespacedName{Name: "test-image", Namespace: "default"}
		originalLister = listVersionsForImage
		stubbedVersions = []string{"1.0.0"}
		listVersionsForImage = func(
			_ logr.Logger,
			_ string,
			_ string,
			_ bool,
			_ helpers.DockerHubAuth,
			_ helpers.AWSPublicECR,
		) []string {
			return append([]string(nil), stubbedVersions...)
		}
	})

	AfterEach(func() {
		listVersionsForImage = originalLister
	})

	createImage := func(mode skopeoiov1alpha1.Mode, status skopeoiov1alpha1.ImageStatus) *skopeoiov1alpha1.Image {
		image := &skopeoiov1alpha1.Image{
			ObjectMeta: metav1.ObjectMeta{
				Name:      imageNamespaced.Name,
				Namespace: imageNamespaced.Namespace,
			},
			Spec: skopeoiov1alpha1.ImageSpec{
				Mode: mode,
				Source: skopeoiov1alpha1.ImageEndpoint{
					ImageName:    "nginx",
					ImageVersion: "1.0.0",
				},
			},
			Status: status,
		}
		Expect(cl.Create(ctx, image)).To(Succeed())
		return image
	}

	createJob := func(version string, complete bool, failed bool) {
		var conditions []batchv1.JobCondition
		switch {
		case failed:
			conditions = []batchv1.JobCondition{{
				Type:   batchv1.JobFailed,
				Status: corev1.ConditionTrue,
			}}
		case complete:
			conditions = []batchv1.JobCondition{{
				Type:   batchv1.JobComplete,
				Status: corev1.ConditionTrue,
			}}
		}

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      helpers.GenerateSkopeoJobName(imageNamespaced.Name, version),
				Namespace: pullJobNamespace(),
			},
			Status: batchv1.JobStatus{Conditions: conditions},
		}
		Expect(cl.Create(ctx, job)).To(Succeed())
	}

	It("marks the image as completed when all jobs succeed", func() {
		image := createImage(skopeoiov1alpha1.ONE_SHOT, skopeoiov1alpha1.ImageStatus{
			History: []skopeoiov1alpha1.History{{PerformedAt: metav1.Now()}},
			Phase:   phaseRunning,
		})
		createJob("1.0.0", true, false)

		updated, err := updateStatusFromJobs(
			reconciler,
			ctx,
			ctrl.Request{NamespacedName: imageNamespaced},
			image,
			logr.Discard(),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		var stored skopeoiov1alpha1.Image
		Expect(cl.Get(ctx, imageNamespaced, &stored)).To(Succeed())
		Expect(stored.Status.Phase).To(Equal(phaseCompleted))
	})

	It("marks the image as failed when a job fails", func() {
		image := createImage(skopeoiov1alpha1.ONE_SHOT, skopeoiov1alpha1.ImageStatus{
			History: []skopeoiov1alpha1.History{{PerformedAt: metav1.Now()}},
			Phase:   phaseRunning,
		})
		createJob("1.0.0", false, true)

		updated, err := updateStatusFromJobs(
			reconciler,
			ctx,
			ctrl.Request{NamespacedName: imageNamespaced},
			image,
			logr.Discard(),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		var stored skopeoiov1alpha1.Image
		Expect(cl.Get(ctx, imageNamespaced, &stored)).To(Succeed())
		Expect(stored.Status.Phase).To(Equal(phaseFailed))
	})

	It("records synced tags only after a job completes in OnceByTag mode", func() {
		stubbedVersions = []string{"1.0.0", "2.0.0"}
		image := createImage(skopeoiov1alpha1.ONCE_BY_TAG, skopeoiov1alpha1.ImageStatus{
			History:          []skopeoiov1alpha1.History{{PerformedAt: metav1.Now()}},
			Phase:            phaseRunning,
			TagAlreadySynced: []string{"2.0.0"},
		})
		createJob("1.0.0", true, false)

		updated, err := updateStatusFromJobs(
			reconciler,
			ctx,
			ctrl.Request{NamespacedName: imageNamespaced},
			image,
			logr.Discard(),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		var stored skopeoiov1alpha1.Image
		Expect(cl.Get(ctx, imageNamespaced, &stored)).To(Succeed())
		Expect(stored.Status.TagAlreadySynced).To(ConsistOf("1.0.0", "2.0.0"))
		Expect(stored.Status.Phase).To(Equal(phaseCompleted))
	})

	It("marks OnceByTag images as completed when every tag is already synced", func() {
		stubbedVersions = []string{}
		image := createImage(skopeoiov1alpha1.ONCE_BY_TAG, skopeoiov1alpha1.ImageStatus{
			History:          []skopeoiov1alpha1.History{{PerformedAt: metav1.Now()}},
			TagAlreadySynced: []string{"1.0.0"},
		})

		updated, err := updateStatusFromJobs(
			reconciler,
			ctx,
			ctrl.Request{NamespacedName: imageNamespaced},
			image,
			logr.Discard(),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())

		var stored skopeoiov1alpha1.Image
		Expect(cl.Get(ctx, imageNamespaced, &stored)).To(Succeed())
		Expect(stored.Status.Phase).To(Equal(phaseCompleted))
	})
})
