import React, { useState } from 'react';
import api from '../api';
import './UrlShortener.css';

const UrlShortener = ({ onUrlCreated }) => {
  const [formData, setFormData] = useState({
    url: '',
    custom_short: '',
    expiry: 24, // hours
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess(null);

    try {
      const dataToSend = { ...formData, expiry: parseInt(formData.expiry, 10) };
      console.log('Sending URL shortening request:', dataToSend);
      const response = await api.shortenUrl(dataToSend);
      console.log('URL shortening response:', response);
      setSuccess(response);
      setFormData({
        url: '',
        custom_short: '',
        expiry: 24,
      });
      if (onUrlCreated) {
        console.log('Calling onUrlCreated with:', response);
        onUrlCreated(response);
      }
    } catch (err) {
      console.error('URL shortening error:', err);
      setError(err.message);
    } finally {
      setLoading(false);
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

  return (
    <div className="url-shortener">
      <div className="shortener-card">
        <h2>Shorten Your URL</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="url">Original URL</label>
            <input
              type="url"
              id="url"
              name="url"
              value={formData.url}
              onChange={handleChange}
              required
              placeholder="https://example.com/very/long/url"
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="custom_short">Custom Short Code (Optional)</label>
              <input
                type="text"
                id="custom_short"
                name="custom_short"
                value={formData.custom_short}
                onChange={handleChange}
                placeholder="my-custom-link"
                maxLength="20"
              />
            </div>

            <div className="form-group">
              <label htmlFor="expiry">Expiration (Hours)</label>
              <select
                id="expiry"
                name="expiry"
                value={formData.expiry}
                onChange={handleChange}
              >
                <option value="1">1 hour</option>
                <option value="6">6 hours</option>
                <option value="24">1 day</option>
                <option value="168">1 week</option>
                <option value="720">1 month</option>
              </select>
            </div>
          </div>

          {error && <div className="error-message">{error}</div>}

          <button type="submit" disabled={loading} className="shorten-button">
            {loading ? 'Shortening...' : 'Shorten URL'}
          </button>
        </form>

        {success && (
          <div className="success-result">
            <h3>URL Shortened Successfully!</h3>
            <div className="result-item">
              <strong>Short URL:</strong>
              <div className="url-display">
                <a href={success.short_url} target="_blank" rel="noopener noreferrer">
                  {success.short_url}
                </a>
                <button
                  onClick={() => copyToClipboard(success.short_url)}
                  className="copy-button"
                >
                  Copy
                </button>
              </div>
            </div>
            <div className="result-item">
              <strong>Original URL:</strong>
              <div className="url-display">
                <span className="original-url">{success.url}</span>
              </div>
            </div>
            <div className="result-meta">
              <span>Expires in: {success.expiry} hours</span>
              <span>Short Code: {success.custom_short}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default UrlShortener;