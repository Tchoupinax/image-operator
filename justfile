default: up

up: 
  export GRAPHQL_API_ENABLED=true && npx nodemon --watch './**/*.go' --signal SIGTERM --exec go run cmd/main.go
