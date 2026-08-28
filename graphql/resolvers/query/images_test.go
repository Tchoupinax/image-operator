package resolvers_test

import (
	"github.com/Tchoupinax/image-operator/graphql/resolvers/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("imageFromUnstructuredObject", func() {
	It("maps phase and spec fields from a Kubernetes object", func() {
		image := resolvers.ImageFromUnstructuredObject(map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":              "demo-image",
				"creationTimestamp": "2024-01-02T03:04:05Z",
			},
			"status": map[string]interface{}{
				"phase": "COMPLETED",
				"history": []interface{}{
					map[string]interface{}{
						"performedAt": "2024-01-03T04:05:06Z",
					},
				},
			},
			"spec": map[string]interface{}{
				"source": map[string]interface{}{
					"name":    "nginx",
					"version": "1.0.0",
				},
				"destination": map[string]interface{}{
					"name":    "registry.example.com/nginx",
					"version": "1.0.0-public",
				},
			},
		})

		Expect(image.Name).To(Equal("demo-image"))
		Expect(image.CreatedAt).To(Equal("2024-01-02T03:04:05Z"))
		Expect(image.Status).To(Equal("COMPLETED"))
		Expect(image.LastExecution).To(Equal("2024-01-03T04:05:06Z"))
		Expect(image.Source.ImageName).To(Equal("nginx"))
		Expect(image.Source.ImageVersion).To(Equal("1.0.0"))
		Expect(image.Destination.ImageName).To(Equal("registry.example.com/nginx"))
		Expect(image.Destination.ImageVersion).To(Equal("1.0.0-public"))
	})

	It("falls back to RUNNING when history exists but phase is empty", func() {
		image := resolvers.ImageFromUnstructuredObject(map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "running-image",
			},
			"status": map[string]interface{}{
				"history": []interface{}{
					map[string]interface{}{
						"performedAt": "2024-01-03T04:05:06Z",
					},
				},
			},
		})

		Expect(image.Status).To(Equal("RUNNING"))
		Expect(image.LastExecution).To(Equal("2024-01-03T04:05:06Z"))
	})
})
