# skopeo-operator Helm Chart

## Credentials

The helm chart creates the secret by default. You can provider your own secret. It should match following properties.

```
stringData:
  credentialsDestinationUsername: "xxx"
  credentialsDestinationPassword: "xxx"
  credentialsSourceUsername: "xxx"
  credentialsSourcePassword: "xxx"
```
