FROM --platform=$BUILDPLATFORM golang:1.16 as builder
ENV TZ=Asia/Shanghai LANG="C.UTF-8"
ARG TARGETARCH
ARG TARGETOS

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

# Copy the go source
COPY queries queries
COPY service service
COPY circle.go circle.go

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o circled circle/cmd/main.go

FROM --platform=$BUILDPLATFORM alpine:latest
WORKDIR /
COPY --from=builder /workspace/circled .
ENTRYPOINT ["/circled"]