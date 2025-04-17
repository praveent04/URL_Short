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

## Environment Setup

For security reasons, this project uses environment variables to store configuration settings. **Do not commit your actual .env file to version control.**

1. Copy the example environment file:
   ```
   cp .env.example .env
   ```

2. Update the `.env` file with your Redis credentials:
   ```
   DOMAIN=localhost:3000
   DB_ADDR=your-redis-host
   DB_PORT=your-redis-port
   DB_PASS=your-redis-password
   DB_USER=default
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

1. Ensure Docker Desktop is running

2. Build and start the containers:
   ```
   docker-compose up -d
   ```

3. The server will be accessible at http://localhost:3000

Note: The Docker Compose file is configured to use environment variables from your `.env` file. Make sure it's properly set up before running.

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

## Docker Compose Troubleshooting

If you encounter issues with Docker Compose:
- Ensure Docker Desktop is running
- Check that your `.env` file exists and has the correct format
- If you see the error `unable to get image` or pipe errors, restart Docker Desktop
- Make sure you're running the Docker commands from the project root directory

## License

MIT
