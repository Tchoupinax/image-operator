package controller

import (
	"context"
	"time"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImageReconciler", func() {
	var (
		ctx            context.Context
		scheme         *runtime.Scheme
		reconciler     *ImageReconciler
		originalLister listVersionsFunc
	)

	BeforeEach(func() {
		ctx = context.Background()
		scheme = runtime.NewScheme()
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

	It("ignores missing images", func() {
		result, err := reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: types.NamespacedName{Name: "missing", Namespace: "default"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))
	})

	It("plans jobs when a new generation is observed", func() {
		image := &skopeoiov1alpha1.Image{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "new-gen",
				Namespace:  "default",
				Generation: 2,
			},
			Spec: skopeoiov1alpha1.ImageSpec{
				Mode: skopeoiov1alpha1.RECURRENT,
				Source: skopeoiov1alpha1.ImageEndpoint{
					ImageName:    "nginx",
					ImageVersion: "1.0.0",
				},
			},
			Status: skopeoiov1alpha1.ImageStatus{
				LastGenerationSeen: 1,
				History: []skopeoiov1alpha1.History{
					{PerformedAt: metav1.NewTime(time.Now().Add(-2 * time.Hour))},
				},
			},
		}
		Expect(reconciler.Create(ctx, image)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: types.NamespacedName{Name: image.Name, Namespace: image.Namespace},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Requeue).To(BeTrue())
		Expect(result.RequeueAfter).To(Equal(60 * time.Second))

		var stored skopeoiov1alpha1.Image
		Expect(reconciler.Get(ctx, types.NamespacedName{Name: image.Name, Namespace: image.Namespace}, &stored)).To(Succeed())
		Expect(stored.Status.LastGenerationSeen).To(Equal(int64(2)))
	})

	It("waits for OneShot images while they are still running", func() {
		image := &skopeoiov1alpha1.Image{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "oneshot-running",
				Namespace:  "default",
				Generation: 1,
			},
			Spec: skopeoiov1alpha1.ImageSpec{
				Mode: skopeoiov1alpha1.ONE_SHOT,
				Source: skopeoiov1alpha1.ImageEndpoint{
					ImageName:    "nginx",
					ImageVersion: "1.0.0",
				},
			},
			Status: skopeoiov1alpha1.ImageStatus{
				LastGenerationSeen: 1,
				Phase:              phaseRunning,
				History: []skopeoiov1alpha1.History{
					{PerformedAt: metav1.Now()},
				},
			},
		}
		Expect(reconciler.Create(ctx, image)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: types.NamespacedName{Name: image.Name, Namespace: image.Namespace},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Requeue).To(BeTrue())
		Expect(result.RequeueAfter).To(Equal(10 * time.Second))
	})

	It("stops reconciling completed OneShot images", func() {
		image := &skopeoiov1alpha1.Image{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "oneshot-done",
				Namespace:  "default",
				Generation: 1,
			},
			Spec: skopeoiov1alpha1.ImageSpec{
				Mode: skopeoiov1alpha1.ONE_SHOT,
				Source: skopeoiov1alpha1.ImageEndpoint{
					ImageName:    "nginx",
					ImageVersion: "1.0.0",
				},
			},
			Status: skopeoiov1alpha1.ImageStatus{
				LastGenerationSeen: 1,
				Phase:              phaseCompleted,
				History: []skopeoiov1alpha1.History{
					{PerformedAt: metav1.Now()},
				},
			},
		}
		Expect(reconciler.Create(ctx, image)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: types.NamespacedName{Name: image.Name, Namespace: image.Namespace},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))
	})
})
