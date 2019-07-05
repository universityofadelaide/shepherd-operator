#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/gopher
COMMIT=$(shell git rev-list -1 HEAD)
VERSION=$(shell git describe --tags --always)
IMAGE=skpr/operator

export KUBECONFIG=${GOPATH}/src/github.com/skpr/local/.kube/config

all: test build

# Run tests
test: generate test-unit test-integration

test-unit:
	go test -tags=unit ./...

test-integration:
	go test -tags=integration ./...

# Build binary
build: generate lint
	GOOS=linux go build -o bin/skpr-operator \
	         -ldflags='-extldflags "-static" -X github.com/skpr/operator/cmd/skpr-operator/version.GitVersion=${VERSION} -X github.com/skpr/operator/cmd/skpr-operator/version.GitCommit=${COMMIT}' \
			 github.com/skpr/operator/cmd/skpr-operator

# Run against the configured Kubernetes cluster in ~/.kube/config
# The following is using database connection details from https://github.com/skpr/local
run: build install
	kubectl apply -f config/local/
	kubedev

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Lint code
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`
	go vet ./pkg/... ./cmd/...
	go fmt ./pkg/... ./cmd/...

# Build the docker image and push
release:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

# https://github.com/ahmetb/gen-crd-api-reference-docs
docs:
	gen-crd-api-reference-docs -config docs/config.json \
	                           -template-dir docs/template \
	                           -api-dir "github.com/skpr/operator/pkg/apis" \
				   -out-file docs/apis.md

.PHONY: *
