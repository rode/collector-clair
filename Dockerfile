# syntax = docker/dockerfile:experimental
# Build the manager binary
FROM golang:1.15-alpine as builder

WORKDIR /workspace

RUN apk add --no-cache gcc libc-dev git

# Copy the Go Modules manifests
COPY go.mod go.sum /workspace/

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY clair clair
COPY listener listener

# Build
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o collector-clair

# Test
RUN go test -v -cover -tags unit ./...

# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as runner
WORKDIR /
COPY --from=builder /workspace/collector-clair .
USER nonroot:nonroot

ENTRYPOINT ["./collector-clair"]
EXPOSE 8080


# - docker build -t collector-clair:buildID --target=builder
# - copy coverage out
#   - docker run collector-clair:buildID
#   - docker cp collector-clair:buildID coverage.txt .
#   - docker rm collector-clair
# - docker build -t collector-clair --target=artifact
