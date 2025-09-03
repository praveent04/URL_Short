# syntax=docker/dockerfile:1

# Stage 1: Build React Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

# Install deps deterministically; requires package-lock.json
COPY frontend/package*.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Build
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# Download modules with cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Stage 3: Final Production Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /srv

# Copy backend binary
COPY --from=backend-builder /app/main /srv/main

# Copy frontend build files
COPY --from=frontend-builder /app/frontend/build /srv/frontend/build

# Optional: run as non-root (create user and switch)
# RUN adduser -D -H app && chown -R app:app /srv
# USER app

EXPOSE 3000
# Do NOT bake environment files into the image; pass env vars at runtime
# ENV STATIC_DIR=/srv/frontend/build
CMD ["/srv/main"]
