![AI generated image showing an octobus manipulating containers](.github/docs/logo.png)

# skopeo-operator

## Install with Helm Chart

```bash
helm repo add skopeo-operator https://tchoupinax.github.io/skopeo-operator
helm repo update

helm upgrade --install skopeo-operator skopeo-operator/skopeo-operator
```

## Usage

According the use-case you have, select the good configuration for you:

> I want to sync an image with a specific tag and I want to do only one

Use mode `OneShot` and fill version with you specific tag (e.g. `v1.2.3`).

> I want to sync a version and I want all new version to be synced automatically

Use mode `OneShot` and fill version with a matching pattern (e.g. `v1.2.x`). It will sync current available version and watch for future version to sync them (`>v1.2.0 & <v1.3.0`)

> I want to copy and refresh an image every day

Use mode `Recurrent` and provide the desired version (e.g. `node:22-alpine`)

### Reccurent task

```yaml
apiVersion: skopeo.io/v1alpha1
kind: Image
metadata:
  name: nginx-alpine
spec:
  frequency: 15m
  mode: Recurrent # Or OneShot
  source:
    name: quay.io/nginx/nginx-ingress
    version: 3.7-alpine
  destination:
    name: tchoupinax/nginx/nginx-ingress
    version: 3.7-alpine
```

### Tags matching pattern

You can order to copy every images matching a pattern. For exemple, if you want to copy every image like `2.13.1`, `2.13.2`, `2.13.3` etc... you can put version as `2.13.x`.
Moreover, if you want to include release candidate you can with the option `allowCandidateRelease: true`

It will create a job for each version detected.

```yaml
apiVersion: skopeo.io/v1alpha1
kind: Image
metadata:
  name: argocd-2-13-rc
spec:
  allowCandidateRelease: true
  mode: OneShot
  source:
    name: quay.io/argoproj/argocd
    version: v2.13.x
  destination:
    name: tchoupinax/argoproj/argocd
    version: v2.13.x
```

### Options

This is an exemple of the resource with full options.

```yaml
apiVersion: skopeo.io/v1alpha1
kind: Image
metadata:
  name: name

spec:
  allowCandidateRelease: false # Activate if you want release which match *.*.*-rc[O-9]+
  mode: OneShot # or Reccurrent
  frequency: "1m" # or [O-9]+h (hour), [O-9]+d (day), [O-9]+w (week)
  source:
    name: source/argoproj/argocd
    version: v2.13.x
  destination:
    name: destination/argoproj/argocd
    version: v2.13.x
```

## Motivation

We aim to use container images exclusively from our internal registry for various reasons. Among these images, a significant portion consists of "base" images that we did not build ourselves. However, the process of copying these base images presents several challenges:

- It may require authentication, adding complexity to the process.
- Copying images for multiple architectures simultaneously can be cumbersome.
- Not everyone in the organization may have the necessary permissions to perform this operation.
- The stability of the process can vary depending on the method used (e.g., CI/CD pipelines).

Among the various open-source projects known for their ability to build and copy images, one stands out for its efficiency in copying images across registries: [Skopeo](https://github.com/containers/skopeo).

How can we industrialize this process? While it’s possible to use CI for this purpose, incorporating a one-off task into the CI pipeline doesn’t seem advantageous. For example, if we only need to copy the nginx:1.1.1-alpine image once, embedding this operation into the CI process isn't relevant or efficient. We need a more suitable approach for handling such single-use cases. The same limitations apply to recurring tasks. While CI systems can be configured to handle such tasks with scheduled jobs or cron-like setups, they are quite limited.

The idea of building an operator around Skopeo originates from the desire to leverage the right tool in conjunction with the operator pattern, to benefit from Kubernetes's simplicity and scalability.

## Description

```mermaid
sequenceDiagram
actor User
User->>Kubernetes: Apply Image resource to Kubernetes
loop Every 5 seconds
    Kubernetes->>Skopeo Operator: Listen resource's events
    alt is one shot job
        Skopeo Operator->>Skopeo Job: Create job
    else is reccurrent job
        alt last execution is old than frequency
            Skopeo Operator->>Skopeo Job: Create job
        else last execution is newer than frequency
            Skopeo Operator->>Skopeo Operator: Do nothing
        end
    end
    Skopeo Job->>Skopeo Job: Copy image accross registries
    Note right of Skopeo Job: Job is performed asyncronaly and<br>has a random duration.<br>Once the pod has finished,<br>it is deleted.
    Skopeo Operator-->>Kubernetes: If the job is recurrent,<br>ask to recall the loop every 5 seconds
end
```

## Configuration

## Helm chart

You can find an exemple of values [here](charts/skopeo-operator/values.yaml).

### Environment variables list

- `CREDS_DESTINATION_PASSWORD`: ""
- `CREDS_DESTINATION_USERNAME`: ""
- `CREDS_SOURCE_PASSWORD`: ""
- `CREDS_SOURCE_USERNAME`: ""
- `DISABLE_DEST_TLS_VERIFICATION`: "false"
- `DISABLE_SRC_TLS_VERIFICATION`: "false"
- `PULL_JOB_NAMESPACE`: "skopeo-operator"
- `SKOPEO_IMAGE`: "quay.io/containers/skopeo"
- `SKOPEO_VERSION`: "v1.16.1"

## Features

- [x] Copy images accross registries
- [x] Copy images recurrently, frequency is configurable
- [x] Authentication on Password and Username
- [x] Basic monitoring and metrics
- [x] Allow to copy release candidates
- [x] Allow to target version following a pattern
  - [x] Quay.io
  - [ ] Dockerhub
  - [ ] AWS public ECR

## Monitoring

Operator exposes a Prometheus route to show basic metrics about operator and how many reload it has been done.

![Show Grafana's graph](.github/docs/metrics.png)

## Development

### Tests

Command to launch a specific test

```bash
go run github.com/onsi/ginkgo/v2/ginkgo -r --randomize-all --randomize-suites --race --trace -cover internal/helpers/
```
