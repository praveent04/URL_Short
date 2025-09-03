#!/usr/bin/env node

// Render.com Environment Setup Helper
// This script helps you set up environment variables in Render

const fs = require('fs');
const path = require('path');

console.log('🚀 Render.com Environment Setup Helper\n');

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

console.log('📋 Environment Variables to set in Render:\n');
console.log('Copy and paste these into your Render service environment variables:\n');

// Important production modifications for Render
const renderOverrides = {
  'PORT': '10000', // Render uses port 10000
  'REACT_APP_API_URL': 'https://your-render-app-url.onrender.com',
  'JWT_SECRET': 'generate-a-strong-secret-for-production',
  'SMTP_USER': 'your-actual-email@gmail.com',
  'SMTP_PASS': 'your-actual-app-password',
  'FROM_EMAIL': 'your-actual-email@gmail.com'
};

// Set Render-specific PORT
console.log(`PORT=10000 # ⚠️  Render requires port 10000`);

Object.entries(envVars).forEach(([key, value]) => {
  if (key === 'PORT') return; // Skip, already handled above
  
  if (renderOverrides[key]) {
    console.log(`${key}=${renderOverrides[key]} # ⚠️  UPDATE THIS VALUE`);
  } else {
    console.log(`${key}=${value}`);
  }
});

console.log('\n🔧 Important Notes for Render:');
console.log('1. Update REACT_APP_API_URL with your actual Render app URL (.onrender.com)');
console.log('2. PORT must be set to 10000 (Render requirement)');
console.log('3. Generate a strong JWT_SECRET for production');
console.log('4. Configure actual email credentials if you want email features');
console.log('5. Free tier services spin down after inactivity');

console.log('\n📖 For detailed setup instructions, see RENDER_SETUP.md');
