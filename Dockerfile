FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/brandingiron ./cmd/brandingiron/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/brandingiron /usr/local/bin/brandingiron
ENV PORT="9040" DATA_DIR="/data"
EXPOSE 9040
HEALTHCHECK --interval=30s --timeout=5s CMD curl -sf http://localhost:9040/health || exit 1
ENTRYPOINT ["brandingiron"]
