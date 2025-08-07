# ~/triple-s/Dockerfile
FROM golang:1.24-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o triple-s .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/triple-s .
COPY --from=builder /app/data /app/data

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./triple-s"]