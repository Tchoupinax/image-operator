package controller_test

import (
	"os"

	"github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	controller "github.com/Tchoupinax/image-operator/internal/controller/skopeo.io"
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
				&v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName:  "repository.source.com",
							UseAwsIRSA: false,
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:    "repository.destination.com",
							ImageVersion: "v4.5.6-public",
							UseAwsIRSA:   false,
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.Template.Spec.Containers[0].Command).To(Equal([]string{"/bin/bash"}))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(Equal([]string{
				"-c",
				"skopeo copy docker://repository.source.com:v4.5.6 docker://repository.destination.com:v4.5.6-public --all --preserve-digests --src-tls-verify=true --dest-tls-verify=true",
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
				&v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName:  "repository.source.com",
							UseAwsIRSA: false,
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:    "repository.destination.com",
							UseAwsIRSA:   true,
							ImageVersion: "v4.5.6-public",
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.Template.Spec.Containers[0].Command).To(Equal([]string{"/bin/bash"}))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(Equal([]string{
				"-c",
				"yum install -y awscli &&\n\t\t\taws ecr get-login-password --region eu-west-1 | skopeo login --username AWS --password-stdin repository.destination.com &&\n\t\t\tskopeo copy docker://repository.source.com:v4.5.6 docker://repository.destination.com:v4.5.6-public --all --preserve-digests --src-tls-verify=true --dest-tls-verify=true",
			}))
		})
	})

	Describe("when using registry credentials", func() {
		It("should pass source and destination credentials to skopeo", func() {
			GinkgoWriter.Println("setting registry credential env vars")
			Expect(os.Setenv("CREDS_DESTINATION_USERNAME", "dest-user")).To(Succeed())
			Expect(os.Setenv("CREDS_DESTINATION_PASSWORD", "dest-pass")).To(Succeed())
			Expect(os.Setenv("CREDS_SOURCE_USERNAME", "src-user")).To(Succeed())
			Expect(os.Setenv("CREDS_SOURCE_PASSWORD", "src-pass")).To(Succeed())
			DeferCleanup(func() {
				os.Unsetenv("CREDS_DESTINATION_USERNAME")
				os.Unsetenv("CREDS_DESTINATION_PASSWORD")
				os.Unsetenv("CREDS_SOURCE_USERNAME")
				os.Unsetenv("CREDS_SOURCE_PASSWORD")
			})

			job := controller.GenerateSkopeoJob(
				&controller.ImageReconciler{},
				nil,
				ctrl.Request{},
				&v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName:  "repository.source.com",
							UseAwsIRSA: false,
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:    "repository.destination.com",
							ImageVersion: "v4.5.6-public",
							UseAwsIRSA:   false,
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.Template.Spec.Containers[0].Args).To(Equal([]string{
				"-c",
				"skopeo copy docker://repository.source.com:v4.5.6 docker://repository.destination.com:v4.5.6-public --all --preserve-digests --dest-creds=dest-user:dest-pass --src-creds=src-user:src-pass --src-tls-verify=true --dest-tls-verify=true",
			}))
		})
	})

	Describe("when configuring job deletion delay", func() {
		It("should use JOB_DELETION_DELAY_SECONDS when it is valid", func() {
			Expect(os.Setenv("JOB_DELETION_DELAY_SECONDS", "42")).To(Succeed())
			DeferCleanup(func() {
				os.Unsetenv("JOB_DELETION_DELAY_SECONDS")
			})

			job := controller.GenerateSkopeoJob(
				&controller.ImageReconciler{},
				nil,
				ctrl.Request{},
				&v1alpha1.Image{
					Spec: v1alpha1.ImageSpec{
						Source: v1alpha1.ImageEndpoint{
							ImageName: "repository.source.com",
						},
						Destination: v1alpha1.ImageEndpoint{
							ImageName:    "repository.destination.com",
							ImageVersion: "v4.5.6-public",
						},
					},
				},
				logr.Logger{},
				"v4.5.6",
			)

			Expect(job.Spec.TTLSecondsAfterFinished).NotTo(BeNil())
			Expect(*job.Spec.TTLSecondsAfterFinished).To(Equal(int32(42)))
		})
	})
})
