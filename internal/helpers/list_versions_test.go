package helpers_test

import (
	"github.com/Tchoupinax/skopeo-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse time", func() {
	It("should parse correctly minutes", func() {
		Expect(helpers.ListVersion("quay.io/nginx/nginx-ingress", "3.x", false)).To(Equal([]string{"3.7"}))
	})

	It("generate regex", func() {
		Expect(helpers.GenerateRegex("2", false)).To(Equal(`2$`))
		Expect(helpers.GenerateRegex("2.x", false)).To(Equal(`2.\d+$`))
		Expect(helpers.GenerateRegex("19.x", false)).To(Equal(`19.\d+$`))
		Expect(helpers.GenerateRegex("2.12.x", false)).To(Equal(`2.12.\d+$`))
		Expect(helpers.GenerateRegex("4.56.x", false)).To(Equal(`4.56.\d+$`))
		Expect(helpers.GenerateRegex("2.12.x", true)).To(Equal(`2.12.\d+(-rc\d)?$`))
		Expect(helpers.GenerateRegex("3.34.x", true)).To(Equal(`3.34.\d+(-rc\d)?$`))
	})
})
