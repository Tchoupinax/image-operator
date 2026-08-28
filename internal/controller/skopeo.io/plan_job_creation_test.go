package controller

import (
	"context"
	"time"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("planJobCreation", func() {
	Describe("addHistory", func() {
		It("keeps at most 20 history entries", func() {
			image := &skopeoiov1alpha1.Image{}
			for range 25 {
				addHistory(image)
			}
			Expect(image.Status.History).To(HaveLen(20))
		})
	})

	Describe("when no versions are available", func() {
		var (
			ctx            context.Context
			reconciler     *ImageReconciler
			originalLister listVersionsFunc
		)

		BeforeEach(func() {
			ctx = context.Background()
			scheme := runtime.NewScheme()
			Expect(skopeoiov1alpha1.AddToScheme(scheme)).To(Succeed())

			reconciler = &ImageReconciler{
				Client: fake.NewClientBuilder().
					WithScheme(scheme).
					WithStatusSubresource(&skopeoiov1alpha1.Image{}).
					Build(),
				Scheme: scheme,
			}

			originalLister = listVersionsForImage
			listVersionsForImage = func(
				_ logr.Logger,
				_ string,
				_ string,
				_ bool,
				_ helpers.DockerHubAuth,
				_ helpers.AWSPublicECR,
			) []string {
				return nil
			}
		})

		AfterEach(func() {
			listVersionsForImage = originalLister
		})

		It("requeues without creating jobs", func() {
			image := &skopeoiov1alpha1.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "empty-image",
					Namespace: "default",
				},
				Spec: skopeoiov1alpha1.ImageSpec{
					Mode: skopeoiov1alpha1.RECURRENT,
					Source: skopeoiov1alpha1.ImageEndpoint{
						ImageName:    "nginx",
						ImageVersion: "missing",
					},
				},
			}
			Expect(reconciler.Create(ctx, image)).To(Succeed())

			result := planJobCreation(
				reconciler,
				ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{Name: image.Name, Namespace: image.Namespace}},
				image,
				logr.Discard(),
			)

			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter).To(Equal(60 * time.Second))
			Expect(image.Status.History).To(BeEmpty())
		})
	})

	Describe("when versions are available", func() {
		var (
			ctx            context.Context
			reconciler     *ImageReconciler
			originalLister listVersionsFunc
		)

		BeforeEach(func() {
			ctx = context.Background()
			scheme := runtime.NewScheme()
			Expect(skopeoiov1alpha1.AddToScheme(scheme)).To(Succeed())
			Expect(batchv1.AddToScheme(scheme)).To(Succeed())

			reconciler = &ImageReconciler{
				Client: fake.NewClientBuilder().
					WithScheme(scheme).
					WithStatusSubresource(&skopeoiov1alpha1.Image{}).
					Build(),
				Scheme: scheme,
			}

			originalLister = listVersionsForImage
			listVersionsForImage = func(
				_ logr.Logger,
				_ string,
				_ string,
				_ bool,
				_ helpers.DockerHubAuth,
				_ helpers.AWSPublicECR,
			) []string {
				return []string{"1.0.0"}
			}
		})

		AfterEach(func() {
			listVersionsForImage = originalLister
		})

		It("sets phase to RUNNING and requeues OneShot images quickly", func() {
			image := &skopeoiov1alpha1.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "oneshot-image",
					Namespace: "default",
				},
				Spec: skopeoiov1alpha1.ImageSpec{
					Mode: skopeoiov1alpha1.ONE_SHOT,
					Source: skopeoiov1alpha1.ImageEndpoint{
						ImageName:    "nginx",
						ImageVersion: "1.0.0",
					},
					Destination: skopeoiov1alpha1.ImageEndpoint{
						ImageName:    "dest.example.com/nginx",
						ImageVersion: "1.0.0",
					},
				},
			}
			Expect(reconciler.Create(ctx, image)).To(Succeed())

			result := planJobCreation(
				reconciler,
				ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{Name: image.Name, Namespace: image.Namespace}},
				image,
				logr.Discard(),
			)

			Expect(result.Requeue).To(BeTrue())
			Expect(result.RequeueAfter).To(Equal(10 * time.Second))
			Expect(image.Status.Phase).To(Equal(phaseRunning))
			Expect(image.Status.History).To(HaveLen(1))
			Expect(image.Status.TagAlreadySynced).To(BeEmpty())
		})
	})
})
