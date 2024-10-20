package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCheckGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helpers suite")
}
