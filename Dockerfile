# Build the manager binary
FROM golang:latest as builder

# Copy in the go src
WORKDIR /go/src/github.com/universityofadelaide/shepherd-operator
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/universityofadelaide/shepherd-operator/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /go/src/github.com/universityofadelaide/shepherd-operator/manager .
ENTRYPOINT ["/manager"]
