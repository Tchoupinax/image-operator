package helpers_test

import (
	"os"

	"github.com/Tchoupinax/image-operator/internal/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetEnv", func() {
	It("returns the environment value when it is set", func() {
		Expect(os.Setenv("TEST_GET_ENV_KEY", "custom")).To(Succeed())
		DeferCleanup(func() {
			os.Unsetenv("TEST_GET_ENV_KEY")
		})

		Expect(helpers.GetEnv("TEST_GET_ENV_KEY", "default")).To(Equal("custom"))
	})

	It("returns the fallback when the environment value is missing", func() {
		Expect(helpers.GetEnv("TEST_GET_ENV_MISSING_KEY", "default")).To(Equal("default"))
	})
})
