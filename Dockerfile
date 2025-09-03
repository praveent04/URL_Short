# syntax=docker/dockerfile:1

# Stage 1: Build Go Backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# Download modules with cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Stage 2: Final Production Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /srv

# Copy backend binary
COPY --from=backend-builder /app/main /srv/main

# Optional: run as non-root (create user and switch)
# RUN adduser -D -H app && chown -R app:app /srv
# USER app

EXPOSE 10000
# Do NOT bake environment files into the image; pass env vars at runtime
CMD ["/srv/main"]
