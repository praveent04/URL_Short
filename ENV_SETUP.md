# Environment Configuration

This project uses a single `.env` file located in the root directory for all environment variables.

## Setup

1. Copy the template:
   ```bash
   cp .env.template .env
   ```

2. Update the values in `.env` with your actual configuration:
   - Database credentials
   - API keys
   - Domain settings
   - Email configuration

## Variables

### Application
- `DOMAIN`: Your domain (e.g., ynit.com)
- `PORT`: Backend server port (default: 3000)

### Frontend
- `FRONTEND_PORT`: Frontend development server port (default: 3001)
- `REACT_APP_API_URL`: API endpoint for frontend

### Database
- `DB_ADDR`: Redis host
- `DB_PORT`: Redis port
- `DB_PASS`: Redis password
- `DB_USER`: Redis username
- `DB_HOST_PG`: PostgreSQL host
- `DB_PORT_PG`: PostgreSQL port
- `DB_USER_PG`: PostgreSQL username
- `DB_PASS_PG`: PostgreSQL password
- `DB_NAME`: PostgreSQL database name

### Email (Optional)
- `SMTP_HOST`: SMTP server host
- `SMTP_PORT`: SMTP server port
- `SMTP_USER`: Email username
- `SMTP_PASS`: Email password (use app password for Gmail)
- `FROM_EMAIL`: Sender email address
- `FROM_NAME`: Sender name

### Security
- `JWT_SECRET`: Secret key for JWT token signing

## Development

To start the development environment:

1. Backend: `go run main.go`
2. Frontend: `cd frontend && npm start`

The frontend will automatically load environment variables from the root `.env` file.

## Production

For production deployment, set environment variables directly in your hosting platform rather than using the `.env` file.
