package helpers_test

import (
	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Array helpers", func() {
	Describe("Filter", func() {
		It("returns elements matching the predicate", func() {
			result := helpers.Filter([]string{"a", "bb", "c", "dd"}, func(s string) bool {
				return len(s) > 1
			})
			Expect(result).To(Equal([]string{"bb", "dd"}))
		})

		It("returns nil for an empty input slice", func() {
			Expect(helpers.Filter(nil, func(string) bool { return true })).To(BeNil())
		})
	})

	Describe("Contains", func() {
		It("returns true when the value is present", func() {
			Expect(helpers.Contains([]string{"v1", "v2"}, "v2")).To(BeTrue())
		})

		It("returns false when the value is absent", func() {
			Expect(helpers.Contains([]string{"v1"}, "v2")).To(BeFalse())
		})
	})
})
