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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	skopeoiov1alpha1 "github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	LastTimeImageWasReloaded prometheus.GaugeVec
	PrometheusReloadGauge    prometheus.CounterVec
	Scheme                   *runtime.Scheme
	client.Client
}

// +kubebuilder:rbac:groups=skopeo.io,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skopeo.io,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=skopeo.io,resources=images/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var image skopeoiov1alpha1.Image
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &image); err != nil {
		logger.Info("Object not found", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	newGeneration := false
	if image.Generation != image.Status.LastGenerationSeen {
		newGeneration = true
		image.Status.LastGenerationSeen = image.Generation
	}

	if len(image.Status.History) == 0 || newGeneration {
		return planJobCreation(r, ctx, req, &image, logger), nil
	}

	if _, err := updateStatusFromJobs(r, ctx, req, &image, logger); err != nil {
		logger.Error(err, "Failed to update image status from jobs", "image", image.Name)
		return ctrl.Result{}, err
	}

	if image.Spec.Mode == skopeoiov1alpha1.ONE_SHOT {
		if image.Status.Phase == phaseRunning {
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		return ctrl.Result{}, nil
	}

	parsedFrequency := helpers.ParseTime(image.Spec.Frequency)

	lastApplication := image.Status.History[len(image.Status.History)-1].PerformedAt
	timeDifference := metav1.Now().Sub(lastApplication.Time)

	result := ctrl.Result{
		Requeue:      true,
		RequeueAfter: 10 * time.Second,
	}

	if timeDifference > parsedFrequency {
		result = planJobCreation(r, ctx, req, &image, logger)
	}

	return result, nil
}

func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skopeoiov1alpha1.Image{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
