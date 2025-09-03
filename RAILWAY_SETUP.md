# Railway Deployment Setup

## Required Environment Variables

Set these environment variables in your Railway project dashboard:

### Application Configuration
```
DOMAIN=ynit.com
PORT=3000
```

### Frontend Configuration
```
FRONTEND_PORT=3001
REACT_APP_API_URL=https://your-railway-app-url.railway.app
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

### Email Configuration
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

1. Go to your Railway project dashboard
2. Click on your service
3. Go to "Variables" tab
4. Add all the environment variables listed above
5. **Important**: Update `REACT_APP_API_URL` to your actual Railway app URL
6. Redeploy the application

## Getting Your Railway App URL

1. Go to your Railway project dashboard
2. Click on your service
3. Go to "Settings" tab
4. Copy the domain from "Public Networking" section
5. Update `REACT_APP_API_URL` environment variable with this URL

## Domain Configuration

If you want to use a custom domain (ynit.com):
1. Go to "Settings" tab in your Railway service
2. Click "Add Custom Domain"
3. Enter your domain name
4. Follow the DNS configuration instructions
5. Update `REACT_APP_API_URL` to use your custom domain
