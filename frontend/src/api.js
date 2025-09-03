const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL;
    this.token = localStorage.getItem('token');
  }

  setToken(token) {
    this.token = token;
    localStorage.setItem('token', token);
  }

  getToken() {
    return this.token || localStorage.getItem('token');
  }

  removeToken() {
    this.token = null;
    localStorage.removeItem('token');
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    // Add authorization header if token exists
    const token = this.getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      console.log('API: Sending request with token:', token.substring(0, 20) + '...');
    } else {
      console.log('API: No token found for request to', endpoint);
    }

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Authentication
  async register(userData) {
    const response = await this.request('/api/v1/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    return response;
  }

  async login(credentials) {
    const response = await this.request('/api/v1/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    if (response.token) {
      this.setToken(response.token);
    }
    return response;
  }

  logout() {
    this.removeToken();
  }

  // URL Management
  async shortenUrl(urlData) {
    return await this.request('/api/v1/shorten', {
      method: 'POST',
      body: JSON.stringify(urlData),
    });
  }

  async getUserUrls() {
    return await this.request('/api/v1/urls');
  }

  async getUrlStats(shortCode) {
    return await this.request(`/api/v1/stats/${shortCode}`);
  }

  // Notifications
  async sendExpirationNotifications() {
    return await this.request('/api/v1/notifications/send', {
      method: 'POST',
    });
  }

  // Health check
  async healthCheck() {
    return await this.request('/api/v1/health');
  }

  // Test JWT token
  async testAuth() {
    return await this.request('/api/v1/test');
  }
}

const apiService = new ApiService();
export default apiService;