FROM golang:1.23.4-alpine AS builder

WORKDIR /bin
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o caching-proxy

FROM scratch

WORKDIR /bin
COPY --from=builder /bin/caching-proxy .
CMD ["./caching-proxy"]
