package graphql

import (
	"fmt"
	"log"
	"net/http"

	mutations "github.com/Tchoupinax/image-operator/graphql/resolvers/mutation"
	query "github.com/Tchoupinax/image-operator/graphql/resolvers/query"
	"github.com/graphql-go/graphql"
	graphl "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var modeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "Mode",
	Values: graphql.EnumValueConfigMap{
		"OneShot": &graphql.EnumValueConfig{
			Value: "OneShot",
		},
		"OnceByTag": &graphql.EnumValueConfig{
			Value: "OnceByTag",
		},
		"Recurrent": &graphql.EnumValueConfig{
			Value: "Recurrent",
		},
	},
})

func StartGraphqlServer() {
	fields := graphl.Fields{
		"images": &graphl.Field{
			Type:    graphql.NewList(query.ImageType),
			Resolve: query.Images,
		},
		"imageBuilders": &graphl.Field{
			Type:    graphql.NewList(query.ImageBuilderType),
			Resolve: query.ImageBuilders,
		},
	}
	mutations := graphl.Fields{
		"createImage": &graphl.Field{
			Type: query.ImageType,
			Args: graphql.FieldConfigArgument{
				"sourceRepositoryName": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"sourceVersion": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"destinationRepositoryName": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"destinationVersion": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"mode": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(modeEnum),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: mutations.CreateImage,
		},
	}
	rootQuery := graphl.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphl.ObjectConfig{Name: "RootMutation", Fields: mutations}
	schemaConfig := graphl.SchemaConfig{Query: graphl.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
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
