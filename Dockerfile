# multi-stage build for nga_grep

# builder stage
FROM golang:1.26-bookworm AS builder
WORKDIR /src

# 使用 build cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o /nga_grep ./main.go

# final stage
FROM debian:bookworm-slim
# create work dir for runtime
WORKDIR /data
# copy binary
COPY --from=builder /nga_grep /nga_grep

# make binary executable
RUN chmod +x /nga_grep

# default entrypoint just lists help
ENTRYPOINT ["/nga_grep"]
