#!/usr/bin/env node

// Railway Environment Setup Helper
// This script helps you set up environment variables in Railway

const fs = require('fs');
const path = require('path');

console.log('🚀 Railway Environment Setup Helper\n');

// Read the .env file
const envPath = path.join(__dirname, '.env');
if (!fs.existsSync(envPath)) {
  console.log('❌ .env file not found. Please create one first.');
  process.exit(1);
}

const envContent = fs.readFileSync(envPath, 'utf8');
const envVars = {};

// Parse .env file
envContent.split('\n').forEach(line => {
  line = line.trim();
  if (line && !line.startsWith('#')) {
    const [key, ...valueParts] = line.split('=');
    if (key && valueParts.length > 0) {
      envVars[key] = valueParts.join('=');
    }
  }
});

console.log('📋 Environment Variables to set in Railway:\n');
console.log('Copy and paste these into your Railway project environment variables:\n');

// Important production modifications
const productionOverrides = {
  'REACT_APP_API_URL': 'https://your-railway-app-url.railway.app',
  'JWT_SECRET': 'generate-a-strong-secret-for-production',
  'SMTP_USER': 'your-actual-email@gmail.com',
  'SMTP_PASS': 'your-actual-app-password',
  'FROM_EMAIL': 'your-actual-email@gmail.com'
};

Object.entries(envVars).forEach(([key, value]) => {
  if (productionOverrides[key]) {
    console.log(`${key}=${productionOverrides[key]} # ⚠️  UPDATE THIS VALUE`);
  } else {
    console.log(`${key}=${value}`);
  }
});

console.log('\n🔧 Important Notes:');
console.log('1. Update REACT_APP_API_URL with your actual Railway app URL');
console.log('2. Generate a strong JWT_SECRET for production');
console.log('3. Configure actual email credentials if you want email features');
console.log('4. All sensitive values are already configured for production use');

console.log('\n📖 For detailed setup instructions, see RAILWAY_SETUP.md');
