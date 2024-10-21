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

package controller

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	skopeoiov1alpha1 "github.com/Tchoupinax/skopeo-operator/api/v1alpha1"
	"github.com/Tchoupinax/skopeo-operator/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	client.Client
	Scheme                *runtime.Scheme
	PrometheusReloadGauge prometheus.CounterVec
}

// +kubebuilder:rbac:groups=skopeo.io,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skopeo.io,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=skopeo.io,resources=images/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Image object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var image skopeoiov1alpha1.Image
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &image); err != nil {
		fmt.Println(fmt.Printf("%s not found (%s)", req.Name, req.Namespace))
		return ctrl.Result{}, nil
	}

	var parsedFrequency time.Duration = 5 * time.Minute
	if image.Spec.Mode == "Recurrent" {
		var parsedFrequencyError error

		parsedFrequency, parsedFrequencyError = helpers.ParseTime(image.Spec.Frequency)
		if parsedFrequencyError != nil {
			logger.Error(parsedFrequencyError, "Error when parsing the frequency")
		}
	}

	if len(image.Status.History) == 0 {
		planJobCreation(r, ctx, req, image, logger)
	} else {
		lastApplication := image.Status.History[len(image.Status.History)-1].PerformedAt
		timeDifference := metav1.Now().Sub(lastApplication.Time)

		if timeDifference > parsedFrequency {
			planJobCreation(r, ctx, req, image, logger)
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 10 * time.Second,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skopeoiov1alpha1.Image{}).
		Complete(r)
}
