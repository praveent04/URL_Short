# Multi-stage build for URL Shortener
# Stage 1: Build React Frontend
FROM node:18-alpine as frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.21-alpine as backend-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 3: Final Production Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/main .

# Copy frontend build files
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# Copy environment file (optional, better to use env vars in production)
COPY .env .

EXPOSE 3000

CMD ["./main"]