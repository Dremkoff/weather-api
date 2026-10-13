# syntax=docker/dockerfile:1
FROM golang:1.27.2-alpine AS builder

WORKDIR /src

# Download dependencies separately to maximize build-cache reuse.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a small Linux binary. This assumes the application does not require CGO.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server \
    cmd/weather-api/main.go

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app -h /app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/server ./server

USER app:app

EXPOSE 8000

ENTRYPOINT ["./server"]
