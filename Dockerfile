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
COPY services services
COPY circle.go circle.go
COPY signal.go signal.go
COPY version.go version.go

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o circled ./cmd/main.go

FROM --platform=$BUILDPLATFORM alpine:latest
# Set China time zone copy from https://cloud.tencent.com/developer/article/1626811
ENV TZ Asia/Shanghai
RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && apk del tzdata \

WORKDIR /
COPY --from=builder /workspace/circled .
ENTRYPOINT ["/circled"]