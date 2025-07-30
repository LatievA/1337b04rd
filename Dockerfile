FROM golang:1.24-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o 1337b04rd .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/1337b04rd .
COPY --from=builder /app/internal/ui/templates /app/internal/ui/templates

EXPOSE 8081
HEALTHCHECK --interval=30s --timeout=10s --retries=5 \
  CMD curl -f http://localhost:8081/health || exit 1

CMD ["./1337b04rd"]