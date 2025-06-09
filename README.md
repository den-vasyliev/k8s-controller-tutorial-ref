# Deployment Informer with client-go

- Added a Go function to start a shared informer for Deployments in the default namespace using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- The function supports both kubeconfig and in-cluster authentication:
  - If inCluster is true, uses in-cluster config.
  - If kubeconfig is set, uses the provided path.
  - One of these must be set; there is no default to `~/.kube/config`.
- Logs add, update, and delete events for Deployments using zerolog.

**Usage:**
```bash
git switch feature/step7-informer
go run main.go --log-level trace --kubeconfig ~/.kube/config server
```
**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file or in-cluster config.
- Watches for Deployment events (add, update, delete) in the `default` namespace and logs them.

---

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.
- `server.go` - fasthttp server
- `Makefile` — Build automation tasks.
- `Dockerfile` — Distroless Dockerfile for secure containerization.
- `.github/workflows/` — GitHub Actions workflows for CI/CD.
- `list.go` - list cli command
- `charts/app` - helm chart
- `pkg/informer` - informer implementation
- `pkg/testutil` - envtest kit

## License

MIT License. See [LICENSE](LICENSE) for details.