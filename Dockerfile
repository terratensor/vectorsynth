# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o vectorsynth ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/vectorsynth .
COPY --from=builder /app/internal/web ./internal/web

RUN mkdir -p /data/vectors

VOLUME /data/vectors

EXPOSE 8080

CMD ["./vectorsynth", "-vectors", "/data/vectors/vectors.txt", "-web", "./internal/web"]