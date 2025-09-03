import React, { useState, useEffect } from 'react';
import Login from './components/Login';
import Register from './components/Register';
import Dashboard from './components/Dashboard';
import api from './api';
import './App.css';

function App() {
  const [user, setUser] = useState(null);
  const [authMode, setAuthMode] = useState('login'); // 'login' or 'register'
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is already logged in
    const token = api.getToken();
    if (token) {
      // First try to get user info from localStorage
      const savedUser = localStorage.getItem('user');
      if (savedUser) {
        try {
          const userData = JSON.parse(savedUser);
          setUser(userData);
          setLoading(false);
          return;
        } catch (error) {
          console.error('Error parsing saved user data:', error);
        }
      }

      // If no saved user data, validate token with backend
      api.testAuth()
        .then((response) => {
          const userData = {
            id: response.user_id,
            email: response.user_email,
            name: response.user_email.split('@')[0] // Use email prefix as name
          };
          setUser(userData);
          localStorage.setItem('user', JSON.stringify(userData));
        })
        .catch((error) => {
          console.error('Token validation failed:', error);
          api.removeToken();
          localStorage.removeItem('user');
          setUser(null);
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      setLoading(false);
    }
  }, []);

  const handleLogin = (userData) => {
    const userInfo = {
      id: userData.id,
      email: userData.email,
      name: userData.name || userData.email.split('@')[0]
    };
    setUser(userInfo);
    localStorage.setItem('user', JSON.stringify(userInfo));
  };

  const handleRegister = (userData) => {
    const userInfo = {
      id: userData.id,
      email: userData.email,
      name: userData.name || userData.email.split('@')[0]
    };
    setUser(userInfo);
    localStorage.setItem('user', JSON.stringify(userInfo));
  };

  const handleLogout = () => {
    api.logout();
    localStorage.removeItem('user');
    setUser(null);
    setAuthMode('login');
  };

  const switchToRegister = () => {
    setAuthMode('register');
  };

  const switchToLogin = () => {
    setAuthMode('login');
  };

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner">Loading...</div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="App">
        {authMode === 'login' ? (
          <Login
            onLogin={handleLogin}
            onSwitchToRegister={switchToRegister}
          />
        ) : (
          <Register
            onRegister={handleRegister}
            onSwitchToLogin={switchToLogin}
          />
        )}
      </div>
    );
  }

  return (
    <div className="App">
      <Dashboard user={user} onLogout={handleLogout} />
    </div>
  );
}

export default App;
