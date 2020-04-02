FROM golang:alpine

RUN apk add --no-cache ca-certificates bash curl make

# Install kubebuilder 1.x
VOLUME /tmp
WORKDIR /tmp
ARG KUBEBUILDER_VERSION=1.0.8
RUN curl -sSL "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz" > kubebuilder.tgz && \
    tar -vxxzf kubebuilder.tgz && \
    mv kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64/bin/* /usr/local/bin/

# Install kustomize 1.x
ARG KUSTOMIZE_VERSION=1.0.11
RUN curl -sSL https://github.com/kubernetes-sigs/kustomize/releases/download/v${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_amd64 > /usr/local/bin/kustomize

RUN chmod +x /usr/local/bin/*

WORKDIR /go