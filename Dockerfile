FROM --platform=$BUILDPLATFORM golang:1.16 as builder
ARG TARGETARCH
ARG TARGETOS

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

# Copy the go source
COPY cmd     cmd
COPY queries queries
COPY service service
COPY circle.go circle.go
COPY signal.go signal.go
COPY config.go config.go

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o circled ./cmd/main.go

FROM --platform=$BUILDPLATFORM alpine:latest
ENV TZ=Asia/Shanghai LANG="C.UTF-8"
WORKDIR /
COPY --from=builder /workspace/circled .
ENTRYPOINT ["/circled"]