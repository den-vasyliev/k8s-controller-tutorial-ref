# app Helm Chart

This Helm chart deploys the app with configurable image repository and tag.

## Usage

Override the image tag to deploy a specific version:

```sh
helm install my-app ./charts/app \
  --set image.repository=ghcr.io/your-org/app \
  --set image.tag=v1.2.3
```

The image tag is set by CI to the Git tag (if present) or the commit SHA.