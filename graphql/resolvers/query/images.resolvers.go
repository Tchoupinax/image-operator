package resolvers

import (
	"context"
	"log"
	"sort"

	"github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	"github.com/graphql-go/graphql"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Image struct {
	AllowCandidateRelease bool
	CreatedAt             string
	Destination           v1alpha1.ImageEndpoint
	Frequency             string
	LastExecution         string
	Mode                  v1alpha1.Mode
	Name                  string
	Source                v1alpha1.ImageEndpoint
	Status                string
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
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"lastExecution": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func ImageFromUnstructuredObject(obj map[string]interface{}) Image {
	var img Image

	if name, found, _ := unstructured.NestedString(obj, "metadata", "name"); found {
		img.Name = name
	}

	if createdAt, found, _ := unstructured.NestedString(obj, "metadata", "creationTimestamp"); found {
		img.CreatedAt = createdAt
	}

	if status, found, _ := unstructured.NestedString(obj, "status", "phase"); found && status != "" {
		img.Status = status
	} else if histories, hasHistory, _ := unstructured.NestedSlice(obj, "status", "history"); hasHistory && len(histories) > 0 {
		img.Status = "RUNNING"
	}

	if histories, found, _ := unstructured.NestedSlice(obj, "status", "history"); found && len(histories) > 0 {
		if performedAt, ok := histories[len(histories)-1].(map[string]interface{})["performedAt"].(string); ok {
			img.LastExecution = performedAt
		}
	}

	if source, found, _ := unstructured.NestedMap(obj, "spec", "source"); found {
		img.Source = v1alpha1.ImageEndpoint{
			ImageName:    source["name"].(string),
			ImageVersion: source["version"].(string),
			UseAwsIRSA:   false,
		}
	}

	if destination, found, _ := unstructured.NestedMap(obj, "spec", "destination"); found {
		img.Destination = v1alpha1.ImageEndpoint{
			ImageName:    destination["name"].(string),
			ImageVersion: destination["version"].(string),
			UseAwsIRSA:   false,
		}
	}

	return img
}

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
		img := ImageFromUnstructuredObject(item.Object)
		images = append(images, img)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].CreatedAt > images[j].CreatedAt
	})

	result := make([]map[string]interface{}, len(images))
	for i, img := range images {
		result[i] = map[string]interface{}{
			"createdAt":     img.CreatedAt,
			"destination":   img.Destination,
			"lastExecution": img.LastExecution,
			"name":          img.Name,
			"source":        img.Source,
			"status":        img.Status,
		}
	}

	return result, nil
}
