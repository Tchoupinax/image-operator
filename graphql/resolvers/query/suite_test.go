package resolvers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraphQLResolvers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GraphQL resolvers suite")
}
