FROM golang:1.12 as builder
WORKDIR /go/src/github.com/skpr/operator
COPY pkg/      pkg/
COPY internal/ internal/
COPY cmd/      cmd/
COPY vendor/   vendor/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/skpr-operator github.com/skpr/operator/cmd/skpr-operator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/skpr/operator/bin/skpr-operator /usr/local/bin/skpr-operator
CMD ["skpr-operator"]
