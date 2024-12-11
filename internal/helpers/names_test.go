package helpers_test

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("when performing names", func() {
	Describe("when generating name for skopeo job", func() {
		It("should truncate to 63 chars if the name is longer", func() {
			Expect(len(helpers.GenerateSkopeoJobName("abcdefghijklmnabcdefghijklmnabcdefghijklmnabcdefghijklmn", "v1.2.3"))).To(Equal(63))
		})
	})
})
