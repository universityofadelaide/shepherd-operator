
# Image URL to use all building/pushing image targets
IMG ?= docker.io/uofa/shepherd-operator:latest
IMG_BUILDER ?= uofa/shepherd-operator:builder-latest
NAMESPACE ?= myproject
SERVICE_ACCOUNT=shepherd

# Disable go modules (use dep)
export GO111MODULE=off

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

ci: ci-test ci-lint
ci-test:
	go test \
	    ./cmd/... \
	    ./pkg/controller/... \
	    ./pkg/utils/...

ci-lint: fmt
	git diff --exit-code

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/universityofadelaide/shepherd-operator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	oc login -u system:admin
	nohup go run ./cmd/manager/main.go --metrics-addr=":8081" &
	sleep 5
	echo "Use tail -f nohup.out to check on backups."
	oc login -u developer -p developer

run-inline: generate fmt vet
	go run ./cmd/manager/main.go --metrics-addr=":8081"

debug: manager
	oc login -u system:admin
	DEBUG="debug" go run ./cmd/manager/main.go --metrics-addr=":8081"

# Install CRDs and RBAC into a cluster
install: manifests
	oc login -u system:admin
	kubectl apply -f config/crds
	kubectl apply -f config/rbac
	oc login -u developer -p developer

kustomize:
	@echo "updating kustomize namespace to ${NAMESPACE}"
	sed -i'' -e 's@namespace: .*@namespace: '"${NAMESPACE}"'@' ./config/default/kustomization.yaml
	docker run --rm -it \
	    -v $(CURDIR):/go/src/github.com/universityofadelaide/shepherd-operator \
	    --workdir /go/src/github.com/universityofadelaide/shepherd-operator \
	    ${IMG_BUILDER} kustomize build config/default -o ./config/deploy.yaml

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests install kustomize
	kubectl apply -f ./config/deploy.yaml

# Generate manifests e.g. CRD, RBAC etc.
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
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml
	docker build -t ${IMG_BUILDER} -f Dockerfile.builder .

# Push the docker image
docker-push:
	docker push ${IMG}
	docker push ${IMG_BUILDER}
