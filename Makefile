APP = k8s-controller-tutorial
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_FLAGS = -v -o $(APP) -ldflags "-X=github.com/den-vasyliev/$(APP)/cmd.appVersion=$(VERSION)"
ENVTEST ?= $(LOCALBIN)/setup-envtest
ENVTEST_VERSION ?= latest
LOCALBIN ?= $(shell pwd)/bin

.PHONY: all build test test-coverage run docker-build clean envtest test-controller

all: build

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
ENVTEST_VERSION ?= release-0.19

format:
	gofmt -s -w ./

lint:
	golint

envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))


build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) main.go

test: envtest
	go install gotest.tools/gotestsum@latest
	CRD_PATH="${PWD}/config/crd/" KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use --bin-dir $(LOCALBIN) -p path)" gotestsum --junitfile report.xml --format testname ./... ${TEST_ARGS}


test-coverage: envtest
	go install github.com/boumenot/gocover-cobertura@latest
	CRD_PATH="${PWD}/config/crd/" KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use --bin-dir $(LOCALBIN) -p path)" go test -coverprofile=coverage.out -covermode=count ./... ${TEST_ARGS}
	go tool cover -func=coverage.out
	gocover-cobertura < coverage.out > coverage.xml

test-controller: envtest
	go install gotest.tools/gotestsum@latest
	CRD_PATH="${PWD}/config/crd/"  KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use --bin-dir $(LOCALBIN) -p path)" gotestsum --junitfile report.xml --format testname ./pkg/ctrl/...

run:
	go run main.go

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(APP):latest .

clean:
	rm -f $(APP)

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef