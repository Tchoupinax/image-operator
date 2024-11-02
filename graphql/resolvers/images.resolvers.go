package resolvers

import (
	"context"
	"log"

	"github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/graphql-go/graphql"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Image struct {
	Name                  string
	AllowCandidateRelease bool
	Destination           v1alpha1.ImageEndpoint
	Frequency             string
	Mode                  v1alpha1.Mode
	Source                v1alpha1.ImageEndpoint
	CreatedAt             string
}

var sourceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Source",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"version": &graphql.Field{
			Type: graphql.String,
		},
		"useAwsIRSA": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var ImageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Image",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"mode": &graphql.Field{
			Type: graphql.String,
		},
		"frequency": &graphql.Field{
			Type: graphql.String,
		},
		"destination": &graphql.Field{
			Type: sourceType,
		},
		"source": &graphql.Field{
			Type: sourceType,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func Images(p graphql.ResolveParams) (interface{}, error) {
	dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "skopeo.io",
		Version:  "v1alpha1",
		Resource: "images",
	}

	customResources, err := dynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing custom resources: %v", err)
	}

	var images []Image
	for _, item := range customResources.Items {
		var img Image

		if name, found, _ := unstructured.NestedString(item.Object, "metadata", "name"); found {
			img.Name = name
		}

		if name, found, _ := unstructured.NestedString(item.Object, "metadata", "creationTimestamp"); found {
			img.CreatedAt = name
		}

		if _, found, _ := unstructured.NestedString(item.Object, "metadata", "name"); found {
			img.Destination = v1alpha1.ImageEndpoint{
				ImageName:    "DD",
				ImageVersion: "de",
				UseAwsIRSA:   false,
			}
			img.Source = v1alpha1.ImageEndpoint{
				ImageName:    "DD",
				ImageVersion: "de",
				UseAwsIRSA:   false,
			}
		}

		images = append(images, img)
	}

	result := make([]map[string]interface{}, len(images))
	for i, img := range images {
		result[i] = map[string]interface{}{
			"destination": img.Destination,
			"name":        img.Name,
			"source":      img.Source,
			"createdAt":   img.CreatedAt,
		}
	}

	return result, nil
}
