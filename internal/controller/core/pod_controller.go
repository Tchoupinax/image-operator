/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var notFoundImageCache []string

type PodReconciler struct {
	client.Client
	Scheme                *runtime.Scheme
	OnFlyNamespaceAllowed []string
}

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var event corev1.Event
	if err := r.Get(ctx, req.NamespacedName, &event); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// Event on the pod
	if event.Reason != "Failed" {
		return ctrl.Result{}, nil
	}

	var shouldThisPodBeHandled = helpers.Contains(r.OnFlyNamespaceAllowed, req.Namespace)
	// If the array equals ["*"], handle every namespace
	if len(r.OnFlyNamespaceAllowed) == 1 && r.OnFlyNamespaceAllowed[0] == "*" {
		shouldThisPodBeHandled = true
	}

	if !shouldThisPodBeHandled {
		logger.Info("Pod " + req.Name + "/" + req.Namespace + " not handled because NS is not allowed")
		return ctrl.Result{}, nil
	}

	// fmt.Println(event.Reason, req.Name)
	// fmt.Println(event.Message)
	// We detect a not found image
	if strings.Contains(event.Message, "not found") {
		data, err := helpers.ExtractImageName(event.Message)
		if err != nil {
			logger.Error(err, "Failed to extract image name")
			return ctrl.Result{}, nil
		}

		var newImage = fmt.Sprintf("%s:%s", data.Image, data.Version)
		// We try to perform the image once only
		if helpers.Contains(notFoundImageCache, newImage) {
			logger.Info("Image detected but already handled: " + newImage)
			return ctrl.Result{}, nil
		}
		notFoundImageCache = append(notFoundImageCache, newImage)

		logger.Info(fmt.Sprintf(
			"Image not found detected => %s/%s:%s — Try to automatically copy it => %s",
			data.Registry, data.Image, data.Version, newImage,
		))

		detectedVersions := helpers.ListVersions(
			logger,
			data.Image,
			data.Version,
			false,
			helpers.DockerHubAuth{},
			helpers.AWSPublicECR{},
		)

		if len(detectedVersions) == 1 {
			var destinationDefaultRegistry = os.Getenv("DESTINATION_DEFAULT_REGISTRY")
			var destinationDefaultIrsaUsage = helpers.GetEnv("DESTINATION_DEFAULT_AWS_IRSA_USAGE", "false")

			if destinationDefaultRegistry == "" {
				logger.Info("On-fly copy aborted because default registry not defined (see DESTINATION_DEFAULT_REGISTRY).")
				return ctrl.Result{}, nil
			}

			var image = v1alpha1.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:      helpers.GenerateImageName(data.Image, detectedVersions[0]),
					Namespace: "image-operator",
				},
				Spec: v1alpha1.ImageSpec{
					Mode: "OneShot",
					Source: v1alpha1.ImageEndpoint{
						ImageName:    data.Image,
						ImageVersion: detectedVersions[0],
					},
					Destination: v1alpha1.ImageEndpoint{
						ImageName:    destinationDefaultRegistry + "/" + data.Image,
						ImageVersion: detectedVersions[0],
						UseAwsIRSA:   destinationDefaultIrsaUsage == "true",
					},
				},
			}

			err := r.Create(ctx, &image)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return builder.ControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Watches(
			&corev1.Event{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				// if !helpers.Contains(allowedNamespaces, object.GetNamespace()) {
				//	return []reconcile.Request{}
				// }

				return []reconcile.Request{{
					NamespacedName: client.ObjectKeyFromObject(object),
				}}
			}),
		).
		Complete(r)
}
