# skopeo-operator

## Install with Helm Chart

```
helm repo add skopeo-operator https://tchoupinax.github.io/skopeo-operator
helm repo update
```

## Configuration

- `CREDS_DESTINATION_PASSWORD`: ""
- `CREDS_DESTINATION_USERNAME`: ""
- `CREDS_SOURCE_PASSWORD`: ""
- `CREDS_SOURCE_USERNAME`: ""
- `PULL_JOB_NAMESPACE`: "skopeo-operator" by default
- `SKOPEO_IMAGE`: "quay.io/containers/skopeo"
- `SKOPEO_VERSION`: "v1.16.1"
