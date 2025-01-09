package helpers_test

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/Tchoupinax/image-operator/internal/helpers"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func IsValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func IsIPv6(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	return ip != nil && strings.Contains(ipAddress, ":")
}

func isValidIp(ipAddress string) bool {
	return IsValidIP4(ipAddress) || IsIPv6(ipAddress)
}

var _ = Describe("DockerHub", func() {
	Describe("when getting remaining quota", func() {
		It("should return data", func() {
			var data = helpers.GetDockerhubLimit(logr.Logger{})

			fmt.Printf("\n~ Remaining is %d.\n", data.Remaining)
			Expect(data.Succeeded).To(BeTrue())
			Expect(isValidIp(data.Ip)).To(BeTrue())
			Expect(data.Limit).To(Equal(100))
			Expect(data.LimitWait).To(Equal(21600))
			Expect(data.RemainingWait).To(Equal(21600))
		})
	})
})
