# Golang Kubernetes Controller Tutorial

This project is a step-by-step tutorial for DevOps and SRE engineers to learn about building Golang CLI applications and Kubernetes controllers. Each step is implemented as a feature branch and includes a README section with explanations and command history.

---

## Running Tests

This project uses [envtest](https://book.kubebuilder.io/reference/envtest.html) and [gotestsum](https://github.com/gotestyourself/gotestsum) for robust, CI-friendly testing.

- To run all tests:
  ```sh
  make test
  ```
  This will automatically download the required Kubernetes API server binaries, set up envtest, and run all tests with gotestsum. A JUnit XML report will be generated as `report.xml`.

- To run tests with coverage:
  ```sh
  make test-coverage
  ```
  This will generate a coverage report as `coverage.out`.

---

## Step 1: Golang CLI Application using Cobra

- Initialized a new CLI application using [cobra-cli](https://github.com/spf13/cobra).
- Provides a basic command-line interface.

**Command history:**
```sh
git checkout -b step1-cobra-cli
cobra-cli init --pkg-name github.com/yourusername/k8s-controller-tutorial
# edited main.go, cmd/root.go
```

---

## Step 2: Zerolog for Log Levels

- Integrated [zerolog](https://github.com/rs/zerolog) for structured logging.
- Supports log levels: info, debug, trace, warn, error.

**Command history:**
```sh
git checkout -b step2-zerolog
go get github.com/rs/zerolog
# edited cmd/root.go to add zerolog logging
```

---

## Step 3: pflag for Log Level Flags

- Added [pflag](https://github.com/spf13/pflag) to support a `--log-level` flag.
- Users can set log level via CLI flag.

**Usage:**
```sh
go run main.go --log-level debug
```

**Command history:**
```sh
git checkout -b step3-pflag-loglevel
# edited cmd/root.go to add log-level flag
```

---

## Step 4: FastHTTP Server Command

- Added a new `server` command using [fasthttp](https://github.com/valyala/fasthttp).
- The command starts a FastHTTP server with a configurable port (default: 8080).
- Uses zerolog for logging.

**Usage:**
```sh
go run main.go server --port 8080
```

**What it does:**
- Starts a FastHTTP server on the specified port.
- Responds with "Hello from FastHTTP!" to any request.

**Command history:**
```sh
git checkout -b step4-fasthttp-server
go get github.com/valyala/fasthttp
# created cmd/server.go, added server command
# added cmd/server_test.go for basic tests
go mod tidy
git add .
git commit -m "step4: add fasthttp server command with port flag"
```

---

## Step 6: List Kubernetes Deployments with client-go

- Added a new `list` command using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- Lists deployments in the default namespace.
- Supports a `--kubeconfig` flag to specify the kubeconfig file for authentication.
- Uses zerolog for error logging.

**Usage:**
```sh
go run main.go list --kubeconfig ~/.kube/config
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file.
- Lists all deployments in the `default` namespace and prints their names.

**Command history:**
```sh
go get k8s.io/client-go@v0.29.0
# created cmd/list.go, added list command
# edited go.mod, go.sum
# run go mod tidy
git add .
git commit -m "step6: add list command for Kubernetes deployments using client-go"
```

---

## Step 7: Deployment Informer with client-go

- Added a new `informer` command using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- Runs a shared informer for Deployments in the default namespace.
- Supports both kubeconfig and in-cluster authentication (flags: `--kubeconfig`, `--in-cluster`).
- Logs add, update, and delete events for Deployments using zerolog.

**Usage:**
```sh
go run main.go informer --kubeconfig ~/.kube/config
go run main.go informer --in-cluster
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file or in-cluster config.
- Watches for Deployment events (add, update, delete) in the `default` namespace and logs them.

**Command history:**
```sh
# created cmd/informer.go, added informer command
git add .
git commit -m "step7: add informer command for Kubernetes deployments using client-go"
```

---

## Step 8: /deployments JSON API Endpoint

- Added a `/deployments` endpoint to the FastHTTP server.
- Returns a JSON array of deployment names from the informer's cache (default namespace).
- Uses the informer's local cache, not a live API call.

**Usage:**
```sh
curl http://localhost:8080/deployments
# Output: ["deployment1","deployment2",...]
```

**What it does:**
- Serves a JSON array of deployment names currently in the informer cache.
- Does not query the Kubernetes API directly for each request (fast, efficient).

**Command history:**
```sh
# updated pkg/informer/informer.go, cmd/server.go
git add .
git commit -m "step8: add /deployments JSON API endpoint to server using informer cache"
```

---

Continue to the next steps for more advanced Kubernetes and controller features! 