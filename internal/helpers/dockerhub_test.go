package helpers_test

import (
	"regexp"
	"strings"

	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

var _ = Describe("DockerHub", func() {
	Describe("when getting remaining quota", func() {
		It("should return data", func() {
			var data = helpers.GetDockerhubLimit(logr.Logger{})

			Expect(data.Succeeded).To(BeTrue())

			Expect(validIP4(data.Ip)).To(BeTrue())
			Expect(data.Limit).To(Equal(100))
			Expect(data.LimitWait).To(Equal(21600))
			Expect(data.RemainingWait).To(Equal(21600))
		})
	})
})
