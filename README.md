# URL Shortener with Analytics

A comprehensive URL shortener service built with Go backend and React frontend, featuring advanced analytics and user management.

## Features

### Backend (Go + Fiber)
- ✅ URL shortening with custom codes
- ✅ Click tracking with detailed analytics
- ✅ Geolocation tracking (IP-based)
- ✅ Device and browser detection
- ✅ User authentication (JWT)
- ✅ PostgreSQL database for persistence
- ✅ Redis caching for performance
- ✅ Email notifications for URL expiration
- ✅ RESTful API

### Frontend (React)
- ✅ Modern React application structure
- ✅ Ready for dashboard implementation
- ✅ Authentication components
- ✅ Analytics visualization (to be implemented)

## Setup

1. **Clone and configure:**
   ```bash
   git clone <repository-url>
   cd url-shortener
   ```

2. **Configure environment variables:**
   - **PostgreSQL**: Already configured for NeonDB cloud database
   - **Redis**: Already configured for Redis Cloud
   - **Email**: Configure SMTP settings for notifications (optional)

3. **Start the application:**
   ```bash
   # Using Docker (recommended)
   docker-compose up --build

   # Or run locally
   go run main.go
   ```

4. **Start the frontend:**
   ```bash
   cd frontend
   npm install
   npm start
   ```

## API Endpoints

### Authentication
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login

### URL Management (Protected)
- `POST /api/v1/shorten` - Create short URL
- `GET /api/v1/urls` - Get user's URLs
- `GET /api/v1/stats/:url` - Get URL analytics

### Notifications (Protected)
- `POST /api/v1/notifications/send` - Send expiration notifications

### Public
- `GET /:url` - Redirect to original URL
- `GET /api/v1/health` - Health check

## Email Notifications

To enable URL expiration notifications:

1. Configure SMTP settings in `.env`:
   ```
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USER=your-email@gmail.com
   SMTP_PASS=your-app-password
   FROM_EMAIL=your-email@gmail.com
   FROM_NAME=URL Shortener
   ```

2. For Gmail, generate an app password at: https://myaccount.google.com/apppasswords

3. Trigger notifications manually:
   ```bash
   curl -X POST http://localhost:3000/api/v1/notifications/send \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## Analytics Features

- **Click Tracking**: Timestamp, IP address, user agent
- **Geolocation**: Country and city detection
- **Device Analytics**: Mobile/desktop/bot detection
- **Browser Detection**: Browser type and version
- **Time-based Analytics**: Daily click aggregation
- **Top Countries**: Most active geographic regions

## Database

The application uses **cloud databases** for production-ready deployment:

- **PostgreSQL**: NeonDB (Serverless PostgreSQL)
- **Redis**: Redis Cloud (Managed Redis service)

### Database Schema

#### Users Table
- User accounts with authentication
- JWT-based secure login

#### URLs Table
- Shortened URLs with metadata
- User associations
- Expiration tracking
- Creation timestamps

#### Clicks Table
- Detailed click analytics
- Geographic data (country/city)
- Device and browser information
- Timestamp tracking
- Referrer information

## Development

### Backend
```bash
go mod tidy
go run main.go
```

### Frontend
```bash
cd frontend
npm install
npm start
```

## Production Deployment

1. Set up PostgreSQL and Redis databases
2. Configure environment variables
3. Build and deploy the Go application
4. Build and deploy the React frontend
5. Set up cron job for automatic expiration notifications

## License

MIT License