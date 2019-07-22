
# Image URL to use all building/pushing image targets
IMG ?= docker.io/uofa/shepherd-operator:latest

# Disable go modules (use dep)
export GO111MODULE=off

all: test manager

preflight:
	# Ensure kubebuilder 1.x is in use.
	kubebuilder version | grep 'KubeBuilderVersion:"1'
	# Ensure kustomzie 1.x is in use.
	kustomize version | grep 'KustomizeVersion:1'

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/universityofadelaide/shepherd-operator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

kustomize:
	@echo "updating kustomize namespace"
	sed -i'' -e 's@namespace: .*@namespace: '"${NAMESPACE}"'@' ./config/default/kustomization.yaml
	kustomize build config/default -o ./config/deploy.yaml

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests install kustomize
	kubectl apply -f ./config/deploy.yaml

# Generate manifests e.g. CRD, RBAC etc.
NAMESPACE=shepherd-dev
SERVICE_ACCOUNT=shepherd
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go rbac --service-account=$(SERVICE_ACCOUNT) --service-account-namespace=$(NAMESPACE)

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate: preflight
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

.PHONY: preflight