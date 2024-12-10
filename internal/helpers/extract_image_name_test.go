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
})
