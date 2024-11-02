package resolvers

import (
	"context"
	"log"

	"github.com/graphql-go/graphql"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ImageBuilder struct {
	Name         string
	Architecture string
	Source       string
	CreatedAt    string
}

var ImageBuilderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ImageBuilder",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"architecture": &graphql.Field{
			Type: graphql.String,
		},
		"source": &graphql.Field{
			Type: sourceType,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func ImageBuilders(p graphql.ResolveParams) (interface{}, error) {
	dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "buildah.io",
		Version:  "v1alpha1",
		Resource: "imagebuilders",
	}

	customResources, err := dynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing custom resources: %v", err)
	}

	var images []ImageBuilder
	for _, item := range customResources.Items {
		var imageBuilder ImageBuilder

		if name, found, _ := unstructured.NestedString(item.Object, "metadata", "name"); found {
			imageBuilder.Name = name
		}

		if name, found, _ := unstructured.NestedString(item.Object, "metadata", "creationTimestamp"); found {
			imageBuilder.CreatedAt = name
		}

		if architecture, found, _ := unstructured.NestedString(item.Object, "spec", "architecture"); found {
			imageBuilder.Architecture = architecture
		}

		images = append(images, imageBuilder)
	}

	result := make([]map[string]interface{}, len(images))
	for i, img := range images {
		result[i] = map[string]interface{}{
			"architecture": img.Architecture,
			"name":         img.Name,
			"createdAt":    img.CreatedAt,
		}
	}

	return result, nil
}
