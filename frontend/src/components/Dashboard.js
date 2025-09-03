import React, { useState, useEffect } from 'react';
import api from '../api';
import UrlShortener from './UrlShortener';
import './Dashboard.css';

const Dashboard = ({ user, onLogout }) => {
  const [urls, setUrls] = useState([]);
  const [selectedUrl, setSelectedUrl] = useState(null);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [statsLoading, setStatsLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('urls');

  useEffect(() => {
    loadUserUrls();
  }, []);

  const loadUserUrls = async () => {
    try {
      const response = await api.getUserUrls();
      setUrls(response.urls || []);
    } catch (error) {
      console.error('Failed to load URLs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUrlCreated = (newUrl) => {
    setUrls([newUrl, ...urls]);
  };

  const loadUrlStats = async (shortCode) => {
    setStatsLoading(true);
    try {
      const response = await api.getUrlStats(shortCode);
      setStats(response);
      setSelectedUrl(shortCode);
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setStatsLoading(false);
    }
  };

  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text);
      alert('Copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return <div className="loading">Loading your dashboard...</div>;
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="header-content">
          <h1>URL Shortener Dashboard</h1>
          <div className="user-info">
            <span>Welcome, {user.name}!</span>
            <button onClick={onLogout} className="logout-button">
              Logout
            </button>
          </div>
        </div>
      </header>

      <div className="dashboard-content">
        <div className="tabs">
          <button
            className={`tab ${activeTab === 'urls' ? 'active' : ''}`}
            onClick={() => setActiveTab('urls')}
          >
            My URLs
          </button>
          <button
            className={`tab ${activeTab === 'create' ? 'active' : ''}`}
            onClick={() => setActiveTab('create')}
          >
            Create Short URL
          </button>
          {stats && (
            <button
              className={`tab ${activeTab === 'stats' ? 'active' : ''}`}
              onClick={() => setActiveTab('stats')}
            >
              Analytics
            </button>
          )}
        </div>

        <div className="tab-content">
          {activeTab === 'create' && (
            <UrlShortener onUrlCreated={handleUrlCreated} />
          )}

          {activeTab === 'urls' && (
            <div className="urls-section">
              <h2>Your Shortened URLs</h2>
              {urls.length === 0 ? (
                <div className="empty-state">
                  <p>You haven't created any short URLs yet.</p>
                  <button
                    onClick={() => setActiveTab('create')}
                    className="create-first-button"
                  >
                    Create Your First URL
                  </button>
                </div>
              ) : (
                <div className="urls-grid">
                  {urls.map((url) => (
                    <div key={url.id} className="url-card">
                      <div className="url-info">
                        <div className="url-short">
                          <strong>Short URL:</strong>
                          <a
                            href={`http://localhost:3000/${url.short_code}`}
                            target="_blank"
                            rel="noopener noreferrer"
                          >
                            localhost:3000/{url.short_code}
                          </a>
                          <button
                            onClick={() => copyToClipboard(`http://localhost:3000/${url.short_code}`)}
                            className="copy-button"
                          >
                            Copy
                          </button>
                        </div>
                        <div className="url-original">
                          <strong>Original:</strong>
                          <span className="original-text">{url.original_url}</span>
                        </div>
                        <div className="url-meta">
                          <span>Created: {formatDate(url.created_at)}</span>
                          <span>Expires: {formatDate(url.expires_at)}</span>
                        </div>
                      </div>
                      <div className="url-actions">
                        <button
                          onClick={() => loadUrlStats(url.short_code)}
                          className="stats-button"
                          disabled={statsLoading}
                        >
                          {statsLoading ? 'Loading...' : 'View Stats'}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {activeTab === 'stats' && stats && (
            <div className="stats-section">
              <h2>Analytics for {stats.url.short_code}</h2>

              <div className="stats-overview">
                <div className="stat-card">
                  <h3>Total Clicks</h3>
                  <div className="stat-value">{stats.stats.total_clicks}</div>
                </div>
                <div className="stat-card">
                  <h3>Created</h3>
                  <div className="stat-value">{formatDate(stats.url.created_at)}</div>
                </div>
                <div className="stat-card">
                  <h3>Expires</h3>
                  <div className="stat-value">{formatDate(stats.url.expires_at)}</div>
                </div>
              </div>

              {stats.stats.clicks_by_date && stats.stats.clicks_by_date.length > 0 && (
                <div className="stats-chart">
                  <h3>Clicks Over Time</h3>
                  <div className="chart-placeholder">
                    {stats.stats.clicks_by_date.map((item, index) => (
                      <div key={index} className="chart-item">
                        <span>{item.date}</span>
                        <div
                          className="chart-bar"
                          style={{ width: `${Math.min(item.count * 20, 300)}px` }}
                        >
                          {item.count}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {stats.stats.top_countries && stats.stats.top_countries.length > 0 && (
                <div className="stats-table">
                  <h3>Top Countries</h3>
                  <table>
                    <thead>
                      <tr>
                        <th>Country</th>
                        <th>Clicks</th>
                      </tr>
                    </thead>
                    <tbody>
                      {stats.stats.top_countries.map((country, index) => (
                        <tr key={index}>
                          <td>{country.country || 'Unknown'}</td>
                          <td>{country.count}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              {stats.stats.device_stats && stats.stats.device_stats.length > 0 && (
                <div className="stats-table">
                  <h3>Device Types</h3>
                  <table>
                    <thead>
                      <tr>
                        <th>Device</th>
                        <th>Clicks</th>
                      </tr>
                    </thead>
                    <tbody>
                      {stats.stats.device_stats.map((device, index) => (
                        <tr key={index}>
                          <td>{device.device_type}</td>
                          <td>{device.count}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;