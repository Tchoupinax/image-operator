package helpers_test

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Extract image name", func() {
	It("should extract data correctly", func() {
		data, err := helpers.ExtractImageName(
			"Failed to pull image \"my.custom.registry.com/subfolder/tchoupinax/image-operator-ui:v3gs\": rpc error:",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Version).To(Equal("v3gs"))
		Expect(data.Registry).To(Equal("my.custom.registry.com/subfolder"))
		Expect(data.Image).To(Equal("tchoupinax/image-operator-ui"))
	})

	It("should extract data correctly 2", func() {
		data, err := helpers.ExtractImageName(
			"Failed to pull image \"aws_account_id.dkr.ecr.region.amazonaws.com/tchoupinax/image-operator-ui:v3gs\": rpc error:",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Version).To(Equal("v3gs"))
		Expect(data.Registry).To(Equal("aws_account_id.dkr.ecr.region.amazonaws.com"))
		Expect(data.Image).To(Equal("tchoupinax/image-operator-ui"))
	})

	It("should handle when latest is implicit (root image)", func() {
		data, err := helpers.ExtractImageName(
			"Failed to pull image \"rg.fr-par.scw.cloud/my-registry/busybox\": rpc error:",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Version).To(Equal("latest"))
		Expect(data.Registry).To(Equal("rg.fr-par.scw.cloud/my-registry"))
		Expect(data.Image).To(Equal("busybox"))
	})

	It("should handle when latest is implicit (standard image)", func() {
		data, err := helpers.ExtractImageName(
			"Failed to pull image \"rg.fr-par.scw.cloud/my-registry/tchoupinax/image-operator\": rpc error:",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Version).To(Equal("latest"))
		Expect(data.Registry).To(Equal("rg.fr-par.scw.cloud/my-registry"))
		Expect(data.Image).To(Equal("tchoupinax/image-operator"))
	})

	It("should handle a complex not found image", func() {
		data, err := helpers.ExtractImageName(
			"rpc error: code = NotFound desc = failed to pull and unpack image \"rg.fr-par.scw.cloud/my-registry/repo/image-operator:test2\": failed to resolve reference \"rg.fr-par.scw.cloud/my-registry/repo/image-operator:test2\": rg.fr-par.scw.cloud/my-registry/repo/image-operator:test2: not found",
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Version).To(Equal("test2"))
		Expect(data.Registry).To(Equal("rg.fr-par.scw.cloud/my-registry"))
		Expect(data.Image).To(Equal("repo/image-operator"))
	})
})
