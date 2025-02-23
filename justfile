default: up

test:
  go run github.com/onsi/ginkgo/v2/ginkgo \
    -r \
    --randomize-all \
    --randomize-suites \
    --trace \
    -cover \
    -coverpkg=./... \
    -coverprofile=coverage.out \
    internal/helpers
  go run github.com/onsi/ginkgo/v2/ginkgo \
    -r \
    --randomize-all \
    --randomize-suites \
    --trace \
    -cover \
    -coverpkg=./... \
    -coverprofile=coverage.out \
    internal/controller/skopeo.io
    
up: 
  export GRAPHQL_API_ENABLED=true && npx nodemon --watch './**/*.go' --signal SIGTERM --exec go run cmd/main.go

docker:
  curl --head -H "Authorization: Bearer $(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)" https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest

docker-go:
  go run github.com/onsi/ginkgo/v2/ginkgo run --focus "when getting remaining quota" internal/helpers

helm: 
  kubebuilder edit --plugins=helm/v1-alpha
