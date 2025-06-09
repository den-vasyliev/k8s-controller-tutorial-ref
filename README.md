# List Kubernetes Deployments with client-go

- Added a new `list` command using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- Lists deployments in the default namespace.
- Supports a `--kubeconfig` flag to specify the kubeconfig file for authentication.
- Uses zerolog for error logging.

**Usage:**
```sh
git switch feature/step6-list-deployments 
go run main.go --log-level debug --kubeconfig ~/.kube/config list
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file.
- Lists all deployments in the `default` namespace and prints their names.

---

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.
- `server.go` - fasthttp server
- `Makefile` — Build automation tasks.
- `Dockerfile` — Distroless Dockerfile for secure containerization.
- `.github/workflows/` — GitHub Actions workflows for CI/CD.
- `list.go` - list cli command

## License

MIT License. See [LICENSE](LICENSE) for details.