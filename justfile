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
    
up: 
  export GRAPHQL_API_ENABLED=true && npx nodemon --watch './**/*.go' --signal SIGTERM --exec go run cmd/main.go
