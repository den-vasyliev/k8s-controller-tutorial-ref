# controller-runtime Deployment Controller

- Integrated [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) into the project.
- Added a deployment controller that logs each reconcile event for Deployments in the default namespace.
- The controller is started alongside the FastHTTP server.

**What it does:**
- Uses controller-runtime's manager to run a controller for Deployments.
- Logs every reconcile event (creation, update, deletion) for Deployments.

**Usage:**
```sh
git switch feature/step9-controller-runtime

go run main.go --log-level trace --kubeconfig  ~/.kube/config server
```

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
- `pkg/ctrl` - controller implementation

## License

MIT License. See [LICENSE](LICENSE) for details.