# multi-stage build for nga_grep

# builder stage — trixie matches the distroless base glibc
FROM golang:1.26-trixie AS builder
WORKDIR /src

# 使用 build cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN git rev-parse --short HEAD > /git_commit.txt

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /nga_grep .

# final stage — distroless base (glibc, OpenSSL, CA certs, tzdata; no shell)
FROM gcr.io/distroless/base-debian13

# 设置时区为北京时间（UTC+8） — distroless 内置 tzdata，直接设 TZ 即可
ENV TZ=Asia/Shanghai

# create work dir
WORKDIR /data

# copy binary and git commit info
COPY --from=builder /nga_grep /nga_grep
COPY --from=builder /git_commit.txt /git_commit.txt

# default entrypoint
ENTRYPOINT ["/nga_grep"]
