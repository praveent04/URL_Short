# URL Shortener

A simple URL shortener service built with Go and Redis.

## Features

- Shorten long URLs
- Custom short URL aliases
- Configurable URL expiration
- Cloud Redis support
- Docker support

## Requirements

- Go 1.19+
- Redis (or Redis Cloud account)
- Docker and Docker Compose (optional)

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DOMAIN=localhost:3000
DB_ADDR=your-redis-host
DB_PORT=your-redis-port
DB_PASS=your-redis-password
DB_USER=your-redis-username
```

## Running Locally

1. Install dependencies:
   ```
   go mod download
   ```

2. Run the application:
   ```
   go run main.go
   ```

3. The server will start on port 3000 (or the port specified in the PORT environment variable).

## Running with Docker

1. Build and start the containers:
   ```
   docker-compose up -d
   ```

2. The server will be accessible at http://localhost:3000.

## API Endpoints

### Shorten URL
- **POST** `/api/v1/shorten`
- Body:
  ```json
  {
    "url": "https://example.com/very/long/url",
    "custom_short": "custom",  // Optional
    "expiry": 48              // Optional, in hours
  }
  ```
- Response:
  ```json
  {
    "url": "https://example.com/very/long/url",
    "short_url": "http://localhost:3000/custom",
    "custom_short": "custom",
    "expiry": 48,
    "rate_limit": 10,
    "rate_reset": 30
  }
  ```

### Access Shortened URL
- **GET** `/:url`
- Redirects to the original URL

### Health Check
- **GET** `/api/v1/health`
- Response:
  ```json
  {
    "status": "ok"
  }
  ```

### Debug URL (for development)
- **GET** `/api/v1/debug/:url`
- Response:
  ```json
  {
    "id": "abc123",
    "original_url": "https://example.com"
  }
  ```

## License

MIT