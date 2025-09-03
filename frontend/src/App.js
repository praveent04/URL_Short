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
      // Validate the token with backend
      api.testAuth()
        .then((response) => {
          setUser({ name: 'User', id: response.user_id, email: response.user_email });
        })
        .catch((error) => {
          console.error('Token validation failed:', error);
          api.removeToken();
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
    setUser(userData);
  };

  const handleRegister = (userData) => {
    setUser(userData);
  };

  const handleLogout = () => {
    api.logout();
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
