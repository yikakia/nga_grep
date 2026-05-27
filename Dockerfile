# multi-stage build for nga_grep

# builder stage
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
    GOOS=linux GOARCH=amd64 go build -o /nga_grep .

# final stage
FROM debian:trixie-slim

# install root CAs for outbound HTTPS requests
# 设置时区为北京时间（UTC+8）
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && dpkg-reconfigure -f noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

# create work dir for runtime
WORKDIR /data
# copy binary
COPY --from=builder /nga_grep /nga_grep
COPY --from=builder /git_commit.txt /git_commit.txt

# make binary executable
RUN chmod +x /nga_grep

# default entrypoint just lists help
ENTRYPOINT ["/nga_grep"]
