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

package buildahio

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	buildahiov1alpha1 "github.com/Tchoupinax/image-operator/api/buildah.io/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
)

// ImageBuilderReconciler reconciles a ImageBuilder object
type ImageBuilderReconciler struct {
	client.Client
	Scheme                  *runtime.Scheme
	ImagebuilderBuildsCount prometheus.Counter
}

// +kubebuilder:rbac:groups=buildah.io,resources=imagebuilders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=buildah.io,resources=imagebuilders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=buildah.io,resources=imagebuilders/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;watch;list
func (r *ImageBuilderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var imageBuilder buildahiov1alpha1.ImageBuilder
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &imageBuilder); err != nil {
		logger.Info(fmt.Sprintf("Object %s not found (%s)", req.Name, req.Namespace))
		return ctrl.Result{}, nil
	}

	var newGeneration = false
	if imageBuilder.Generation != imageBuilder.Status.LastGenerationSeen {
		newGeneration = true
		imageBuilder.Status.LastGenerationSeen = imageBuilder.Generation
		updateError := r.Status().Update(ctx, &imageBuilder)
		if updateError != nil {
			fmt.Println(updateError)
		}
	}

	// If it's the first time, start a job creation
	// Also, if the object has changed we reapply it
	if len(imageBuilder.Status.RanJobs) == 0 || newGeneration {
		return CreateBuildahJobs(
			r,
			ctx,
			req,
			imageBuilder,
			logger,
		)
	} else {
		var jobsStatus []string

		for _, jobName := range imageBuilder.Status.RanJobs {
			var job batchv1.Job
			namespacedName := types.NamespacedName{
				Namespace: "image-operator",
				Name:      jobName,
			}
			getJobError := r.Client.Get(ctx, namespacedName, &job)
			if getJobError != nil {
				fmt.Println(getJobError)
			}

			isCompleted := false
			for _, condition := range job.Status.Conditions {
				if condition.Type == batchv1.JobComplete && condition.Status == v1.ConditionTrue {
					isCompleted = true
					break
				}
			}

			if isCompleted {
				jobsStatus = append(jobsStatus, "Completed")
			} else {
				jobsStatus = append(jobsStatus, "NotCompleted")
			}
		}

		if allValuesCompleted(jobsStatus) && !imageBuilder.Status.Pushed {
			CreateBuildahManifestJobs(
				r,
				ctx,
				req,
				imageBuilder,
				logger,
			)

			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

// https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html#setup
// https://yash-kukreja-98.medium.com/develop-on-kubernetes-series-demystifying-the-for-vs-owns-vs-watches-controller-builders-in-c11ab32a046e

// SetupWithManager sets up the controller with the Manager.
func (r *ImageBuilderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&buildahiov1alpha1.ImageBuilder{}).
		Owns(&batchv1.Job{}).
		// TODO: investigate what is this option? => OnlyMetadata
		// Owns(&batchv1.Job{}, builder.OnlyMetadata).
		Complete(r)
}

func allValuesCompleted(arr []string) bool {
	for _, v := range arr {
		if v != "Completed" {
			return false
		}
	}
	return true
}
