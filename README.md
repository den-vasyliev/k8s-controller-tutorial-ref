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

## Informer Test Coverage

The file `pkg/informer/informer_test.go` contains three main test functions:

1. **TestStartDeploymentInformer**
   - Tests the deployment informer event handling and ensures deployment add events are captured.
2. **TestGetDeploymentName**
   - Unit test for the `getDeploymentName` utility, checking both valid and invalid input cases.
3. **TestStartDeploymentInformer_CoversFunction**
   - Ensures the `StartDeploymentInformer` function runs without error.

Each test runs independently when executing `go test ./pkg/informer`. This provides coverage for both informer event handling and utility logic.


## Testing with envtest and Inspecting with kubectl

This project uses [envtest](https://book.kubebuilder.io/reference/envtest.html) to spin up a local Kubernetes API server for integration tests. The test environment writes a kubeconfig to `/tmp/envtest.kubeconfig` so you can inspect the in-memory cluster with `kubectl` while tests are running.

### How to Run and Inspect

1. **Run the informer test:**
   ```sh
   go test ./pkg/informer -run TestStartDeploymentInformer
   ```
   This will:
   - Start envtest and create sample Deployments
   - Write a kubeconfig to `/tmp/envtest.kubeconfig`
   - Sleep for 5 minutes at the end of the test so you can inspect the cluster

2. **In another terminal, use kubectl:**
   ```sh
   kubectl --kubeconfig=/tmp/envtest.kubeconfig get all -A
   kubectl --kubeconfig=/tmp/envtest.kubeconfig get deployments -n default
   kubectl --kubeconfig=/tmp/envtest.kubeconfig describe pod -n default
   ```
   You can use any standard kubectl commands to inspect resources created by the test.

3. **Notes:**
   - The envtest cluster only exists while the test is running. Once the test finishes, the API server is shut down and the kubeconfig is no longer valid.
   - You can adjust the sleep duration in `TestStartDeploymentInformer` if you need more or less time for inspection.

---

For more details, see the code in `pkg/testutil/envtest.go` and `pkg/informer/informer_test.go`.

## License

MIT License. See [LICENSE](LICENSE) for details.