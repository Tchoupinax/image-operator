package helpers_test

import (
	"time"

	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse time", func() {
	It("should parse correctly minutes", func() {
		Expect(helpers.ParseTime("3m")).To(Equal(3 * time.Minute))
	})

	It("should parse correctly hours", func() {
		Expect(helpers.ParseTime("6h")).To(Equal(6 * time.Hour))
	})

	Describe("should parse correctly days", func() {
		It("3d", func() {
			Expect(helpers.ParseTime("3d")).To(Equal(3 * 24 * time.Hour))
		})
		It("5d", func() {
			Expect(helpers.ParseTime("5d")).To(Equal(5 * 24 * time.Hour))
		})
		It("23d", func() {
			Expect(helpers.ParseTime("23d")).To(Equal(23 * 24 * time.Hour))
		})
	})

	Describe("should parse correctly weeks", func() {
		It("1w", func() {
			Expect(helpers.ParseTime("1w")).To(Equal(7 * 24 * time.Hour))
		})
		It("3w", func() {
			Expect(helpers.ParseTime("3w")).To(Equal(21 * 24 * time.Hour))
		})
	})
})
