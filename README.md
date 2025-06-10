# Platform API (CRUD + Swagger)

- Added RESTful CRUD API endpoints for the FrontendPage CRD using FastHTTP and fasthttprouter.
- API handlers use the controller-runtime client to create, update, delete, and list FrontendPage resources in Kubernetes, triggering reconciliation.
- Integrated [Swagger](https://swagger.io/) documentation and served Swagger UI for easy API exploration.
- All API endpoints are under `/api/frontendpages`.

**Usage:**
```sh
git switch feature/step12-platform-api 
go run main.go --log-level trace --kubeconfig  ~/.kube/config server

curl -X POST http://localhost:8080/api/frontendpages -H 'Content-Type: application/json' -d '{"metadata":{"name":"my-page"},"spec":{"contents":"<h1>Hello</h1>","image":"nginx:latest","replicas":2}}'
curl http://localhost:8080/api/frontendpages
curl http://localhost:8080/api/frontendpages/my-page
curl -X PUT http://localhost:8080/api/frontendpages/my-page -H 'Content-Type: application/json' -d '{"spec":{"contents":"<h1>Updated</h1>","image":"nginx:alpine","replicas":1}}'
curl -X DELETE http://localhost:8080/api/frontendpages/my-page
```
- Visit `http://localhost:8080/swagger/index.html` for interactive API docs.

**What it does:**
- Exposes CRUD API for FrontendPage resources, backed by Kubernetes CRDs and controller logic.
- Provides OpenAPI/Swagger docs and UI for easy testing and documentation.

## Running Tests

This project uses [envtest](https://book.kubebuilder.io/reference/envtest.html) and controller-runtime for integration and controller tests.

### Prerequisites
- Go (see go.mod for version)
- Make
- The `setup-envtest` binary (automatically handled by the Makefile)
- CRD YAMLs present in `config/crd/`

### Run all tests
```sh
make test
```
This will:
- Download and set up envtest if needed
- Run all Go tests in the project (including controller and utility tests)

### Run only controller tests
```sh
make test-controller
```
This will:
- Run only the tests in `pkg/ctrl/` (controller logic)

### Test output
- Test logs will show simulated etcd state and resource changes for CRDs and controllers.
- JUnit XML and coverage reports are generated as `report.xml` and `coverage.xml`.

### Troubleshooting
- If you see errors about missing CRDs, ensure you have generated CRDs in `config/crd/` (see Usage above for controller-gen command).
- If you see errors about envtest, try running `make envtest` to ensure the binary is present in `bin/`.
- If you add new CRDs or controllers, re-run `controller-gen` and re-run tests.

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
- `config/crd` - CRD definition
- `pkg/apis` - CRD types and deepcopy
- `pkg/api` - API for PE integration

## License

MIT License. See [LICENSE](LICENSE) for details.