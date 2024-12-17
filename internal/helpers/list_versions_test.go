package helpers_test

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List version from external registry", func() {
	Describe("when image is in Quay.io registry", func() {
		It("should correctly find image for 3.x", func() {
			var expectedValue = []string{
				"3.0",
				"3.1",
				"3.2",
				"3.3",
				"3.4",
				"3.5",
				"3.6",
				"3.7",
				"v3.6",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "quay.io/nginx/nginx-ingress", "3.x", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})

		It("should correctly find image for 3.3.x", func() {
			var expectedValue = []string{
				"3.3.0",
				"3.3.1",
				"3.3.2",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "quay.io/nginx/nginx-ingress", "3.3.x", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})

		It("should correctly find image with release candidate", func() {
			var expectedValue = []string{
				"v2.13.0", "v2.13.0-rc1", "v2.13.0-rc2", "v2.13.0-rc3", "v2.13.0-rc4", "v2.13.0-rc5", "v2.13.1", "v2.13.2",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "quay.io/argoproj/argocd", "2.13.x", true, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})
	})

	Describe("when image is in Dockerhub registry", func() {
		Describe("when image is not prefixed because it is from Dockerhub", func() {
			It("should correctly find image for 2.x", func() {
				var expectedValue = []string{"3.20"}
				Expect(helpers.ListVersions(logr.Logger{}, "alpine", "3.20", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
			})
		})

		It("should correctly find image for 2.x", func() {
			var expectedValue = []string{
				"2.1", "2.10", "2.11", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9", "v2.1", "v2.10", "v2.11", "v2.2", "v2.3", "v2.4", "v2.5", "v2.6", "v2.7", "v2.8", "v2.9",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "library/traefik", "2.x", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})
	})

	Describe("when image is in AWS Public registry", func() {
		It("should correctly find image Node.js 20", func() {
			var expectedValue = []string{
				"20",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "public.ecr.aws/docker/library/node", "20", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})

		It("should correctly find image Node.js 22", func() {
			var expectedValue = []string{
				"22",
			}
			Expect(helpers.ListVersions(logr.Logger{}, "public.ecr.aws/docker/library/node", "22", false, helpers.DockerHubAuth{}, helpers.AWSPublicECR{})).To(Equal(expectedValue))
		})
	})
})

var _ = Describe("when computing regex", func() {
	It("should generate regex", func() {
		Expect(helpers.GenerateRegex("19.x", false)).To(Equal(`^v?19.\d+$`))
		Expect(helpers.GenerateRegex("2", false)).To(Equal(`^v?2$`))
		Expect(helpers.GenerateRegex("2.12.x", false)).To(Equal(`^v?2.12.\d+$`))
		Expect(helpers.GenerateRegex("2.12.x", true)).To(Equal(`^v?2.12.\d+(-rc\d)?$`))
		Expect(helpers.GenerateRegex("2.x", false)).To(Equal(`^v?2.\d+$`))
		Expect(helpers.GenerateRegex("3.34.x", true)).To(Equal(`^v?3.34.\d+(-rc\d)?$`))
		Expect(helpers.GenerateRegex("3.x.x", true)).To(Equal(`^v?3.\d+.\d+(-rc\d)?$`))
		Expect(helpers.GenerateRegex("3.x.0", true)).To(Equal(`^v?3.\d+.0(-rc\d)?$`))
		Expect(helpers.GenerateRegex("4.56.x", false)).To(Equal(`^v?4.56.\d+$`))
		Expect(helpers.GenerateRegex("4.56.x", false)).To(Equal(`^v?4.56.\d+$`))
	})
})
