package controller_test

import (
	"github.com/Tchoupinax/skopeo-operator/api/v1alpha1"
	"github.com/Tchoupinax/skopeo-operator/internal/controller"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Create Skopeo job", func() {
	Describe("when using classic way", func() {
		It("should correctly call the command with given parameters", func() {
			job := controller.GenerateSkopeoJob(
				&controller.ImageReconciler{
					PrometheusReloadGauge: *prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "skopeo_operator_reload_total",
							Help: "Number of reload proccessed",
						},
						[]string{"image"},
					),
				},
				nil,
				ctrl.Request{},
				v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName:  "repository.source.com",
							UseAwsIRSA: false,
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:  "repository.destination.com",
							UseAwsIRSA: false,
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.Template.Spec.Containers[0].Command).To(Equal([]string{"/bin/bash"}))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(Equal([]string{
				"-c",
				"skopeo copy docker://repository.source.com:v4.5.6 docker://repository.destination.com:v4.5.6 --all --preserve-digests --src-tls-verify=true --dest-tls-verify=true",
			}))
		})
	})

	Describe("when destination or source uses AWS IRSA", func() {
		It("should login to registry with AWS CLI", func() {
			job := controller.GenerateSkopeoJob(
				&controller.ImageReconciler{
					PrometheusReloadGauge: *prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "skopeo_operator_reload_total",
							Help: "Number of reload proccessed",
						},
						[]string{"image"},
					),
				},
				nil,
				ctrl.Request{},
				v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName:  "repository.source.com",
							UseAwsIRSA: false,
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:  "repository.destination.com",
							UseAwsIRSA: true,
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.Template.Spec.Containers[0].Command).To(Equal([]string{"/bin/bash"}))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(Equal([]string{
				"-c",
				"yum install -y awscli &&\n\t\t\taws ecr get-login-password --region eu-west-1 | skopeo login --username AWS --password-stdin 326954429656.dkr.ecr.eu-west-1.amazonaws.com &&\n\t\t\tskopeo copy docker://repository.source.com:v4.5.6 docker://repository.destination.com:v4.5.6 --all --preserve-digests --src-tls-verify=true --dest-tls-verify=true",
			}))
		})
	})
})
