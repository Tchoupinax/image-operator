package graphql

import (
	"fmt"
	"log"
	"net/http"

	resolvers "github.com/Tchoupinax/image-operator/graphql/resolvers"
	"github.com/graphql-go/graphql"
	graphl "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func StartGraphqlServer() {
	// Schema
	fields := graphl.Fields{
		"images": &graphl.Field{
			Type:    graphql.NewList(resolvers.ImageType),
			Resolve: resolvers.Images,
		},
		"imageBuilders": &graphl.Field{
			Type:    graphql.NewList(resolvers.ImageBuilderType),
			Resolve: resolvers.ImageBuilders,
		},
	}
	rootQuery := graphl.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphl.SchemaConfig{Query: graphl.NewObject(rootQuery)}
	schema, err := graphl.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// create a graphl-go HTTP handler with our previously defined schema
	// and we also set it to return pretty JSON output
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", cors(h))
	go func() {
		fmt.Println("GraphQL server started")
		httpServerError := http.ListenAndServe(":9090", nil)
		if err := httpServerError; err != nil {
			fmt.Println(err)
		}
	}()
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
