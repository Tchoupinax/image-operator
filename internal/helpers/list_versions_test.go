package helpers_test

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterMatchingTags", func() {
	It("matches Quay nginx 3.x tags", func() {
		tags := []string{
			"2.9", "3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "v3.6", "4.0",
		}
		Expect(helpers.FilterMatchingTags(tags, "3.x", false)).To(Equal([]string{
			"3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "v3.6",
		}))
	})

	It("matches Quay nginx 3.3.x tags", func() {
		tags := []string{"3.3.0", "3.3.1", "3.3.2", "3.4.0"}
		Expect(helpers.FilterMatchingTags(tags, "3.3.x", false)).To(Equal([]string{
			"3.3.0", "3.3.1", "3.3.2",
		}))
	})

	It("matches Quay release candidate tags when enabled", func() {
		tags := []string{
			"v2.13.0", "v2.13.0-rc1", "v2.13.0-rc2", "v2.13.0-rc3", "v2.13.0-rc4",
			"v2.13.0-rc5", "v2.13.1", "v2.13.2", "v2.13.3", "v2.14.0",
		}
		Expect(helpers.FilterMatchingTags(tags, "2.13.x", true)).To(Equal([]string{
			"v2.13.0", "v2.13.0-rc1", "v2.13.0-rc2", "v2.13.0-rc3", "v2.13.0-rc4",
			"v2.13.0-rc5", "v2.13.1", "v2.13.2", "v2.13.3",
		}))
	})

	It("matches Docker Hub alpine exact tags", func() {
		tags := []string{"3.19", "3.20", "3.20.1", "3.21"}
		Expect(helpers.FilterMatchingTags(tags, "3.20", false)).To(Equal([]string{"3.20"}))
	})

	It("matches Docker Hub traefik 2.x tags", func() {
		tags := []string{
			"2.10", "2.11", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9",
			"v2.10", "v2.11", "v2.2", "v2.3", "v2.4", "v2.5", "v2.6", "v2.7", "v2.8", "v2.9",
			"3.0",
		}
		Expect(helpers.FilterMatchingTags(tags, "2.x", false)).To(Equal([]string{
			"2.10", "2.11", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9",
			"v2.10", "v2.11", "v2.2", "v2.3", "v2.4", "v2.5", "v2.6", "v2.7", "v2.8", "v2.9",
		}))
	})

	It("matches AWS Public ECR Node.js major tags", func() {
		tags := []string{"18", "20", "20-alpine", "22", "22-bookworm"}
		Expect(helpers.FilterMatchingTags(tags, "20", false)).To(Equal([]string{"20"}))
		Expect(helpers.FilterMatchingTags(tags, "22", false)).To(Equal([]string{"22"}))
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