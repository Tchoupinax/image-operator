package mutations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Tchoupinax/image-operator/api/skopeo.io/v1alpha1"
	resolvers "github.com/Tchoupinax/image-operator/graphql/resolvers/query"
	"github.com/graphql-go/graphql"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateImage(params graphql.ResolveParams) (interface{}, error) {
	dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())

	fmt.Println(params.Args)

	gvr := schema.GroupVersionResource{
		Group:    "skopeo.io",
		Version:  "v1alpha1",
		Resource: "images",
	}

	image := v1alpha1.Image{
		TypeMeta: v1.TypeMeta{
			APIVersion: "skopeo.io/v1alpha1",
			Kind:       "Image",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      params.Args["name"].(string),
			Namespace: "image-operator",
		},
		Spec: v1alpha1.ImageSpec{
			Mode: "OneShot",
			Destination: v1alpha1.ImageEndpoint{
				ImageName:    params.Args["destinationRepositoryName"].(string),
				ImageVersion: params.Args["destinationVersion"].(string),
			},
			Source: v1alpha1.ImageEndpoint{
				ImageName:    params.Args["sourceRepositoryName"].(string),
				ImageVersion: params.Args["sourceVersion"].(string),
			},
		},
	}

	imageMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&image)
	if err != nil {
		log.Fatalf("Failed to convert Image to unstructured: %v", err)
	}
	unstructuredImage := &unstructured.Unstructured{Object: imageMap}

	_, err = dynamicClient.Resource(gvr).Namespace("image-operator").Create(context.Background(), unstructuredImage, v1.CreateOptions{})
	if err != nil {
		fmt.Printf("Fail to create image: %v", err)

		if strings.Contains(err.Error(), "already exists") {
			imageFromCluster, _ := dynamicClient.Resource(gvr).Namespace("image-operator").Get(context.Background(), image.Name, v1.GetOptions{})

			var localImage resolvers.Image
			if name, found, _ := unstructured.NestedString(imageFromCluster.Object, "metadata", "name"); found {
				localImage.Name = name
			}
			if createdAt, found, _ := unstructured.NestedString(imageFromCluster.Object, "metadata", "creationTimestamp"); found {
				localImage.CreatedAt = createdAt
			}

			return map[string]interface{}{
				"name":      localImage.Name,
				"createdAt": localImage.CreatedAt,
			}, nil
		}
	}

	return map[string]interface{}{
		"name":      image.Name,
		"createdAt": image.ObjectMeta.CreationTimestamp,
	}, nil
}
