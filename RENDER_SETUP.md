# Render.com Deployment Setup

## Required Environment Variables

Set these environment variables in your Render service dashboard:

### Application Configuration
```
DOMAIN=ynit.com
PORT=10000
```
**Note**: Render uses port 10000 by default for web services.

### Frontend Configuration
```
FRONTEND_PORT=3001
REACT_APP_API_URL=https://your-render-app-url.onrender.com
```

### Redis Configuration
```
DB_ADDR=redis-15215.c264.ap-south-1-1.ec2.redns.redis-cloud.com
DB_PORT=15215
DB_PASS=sYlrE16U34UgptoDaSF64lnzY6JNAoSD
DB_USER=default
```

### PostgreSQL Configuration
```
DB_HOST_PG=ep-silent-sea-a19s4xcm-pooler.ap-southeast-1.aws.neon.tech
DB_PORT_PG=5432
DB_USER_PG=neondb_owner
DB_PASS_PG=npg_uJnBPG6XapN1
DB_NAME=neondb
```

### Email Configuration (Optional)
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=your-email@gmail.com
FROM_NAME=YNIT URL Shortener
```

### Security
```
JWT_SECRET=your-super-secret-jwt-key-here-change-this-in-production
```

## Setup Steps

1. Go to your Render dashboard
2. Click on your web service
3. Go to "Environment" tab
4. Add all the environment variables listed above
5. **Important**: Update `REACT_APP_API_URL` to your actual Render app URL
6. **Important**: Set `PORT=10000` (Render's default)
7. Trigger a manual deploy

## Getting Your Render App URL

1. Go to your Render dashboard
2. Click on your service
3. Your URL will be shown at the top (format: `https://your-app-name.onrender.com`)
4. Update `REACT_APP_API_URL` environment variable with this URL

## Build Configuration

Make sure your Render service is configured as:
- **Environment**: Docker
- **Build Command**: (automatic from Dockerfile)
- **Start Command**: `./main`

## Domain Configuration

If you want to use a custom domain (ynit.com):
1. Go to "Settings" tab in your Render service
2. Scroll down to "Custom Domains"
3. Click "Add Custom Domain"
4. Enter your domain name
5. Follow the DNS configuration instructions
6. Update `REACT_APP_API_URL` to use your custom domain

## Render-Specific Notes

- Render automatically sets the `PORT` environment variable to 10000
- Your app must listen on `0.0.0.0:$PORT` (which our Go app already does)
- Environment variables are set in the Render dashboard, not in files
- Free tier services spin down after inactivity (takes ~30 seconds to spin back up)
