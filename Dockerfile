# multi-stage build for nga_grep

# builder stage
FROM golang:1.21-alpine AS builder
WORKDIR /src

# copy go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# copy entire source
COPY . .

# build binary
RUN GOOS=linux GOARCH=amd64 go build -o /nga_grep ./main.go

# final stage
FROM alpine:3.18
# create work dir for runtime
WORKDIR /data
# copy binary
COPY --from=builder /nga_grep /nga_grep

# make binary executable
RUN chmod +x /nga_grep

# default entrypoint just lists help
ENTRYPOINT ["/nga_grep"]
