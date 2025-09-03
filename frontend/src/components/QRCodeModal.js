import React, { useRef } from 'react';
import { QRCodeCanvas } from 'qrcode.react';
import './QRCodeModal.css';

const QRCodeModal = ({ url, shortCode, onClose }) => {
  const qrRef = useRef();

  const downloadQR = () => {
    const canvas = qrRef.current.querySelector('canvas');
    if (canvas) {
      const link = document.createElement('a');
      link.download = `qrcode-${shortCode}.png`;
      link.href = canvas.toDataURL();
      link.click();
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(url);
      alert('URL copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  return (
    <div className="qr-modal-overlay" onClick={onClose}>
      <div className="qr-modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="qr-modal-header">
          <h3>QR Code for {shortCode}</h3>
          <button className="qr-close-button" onClick={onClose}>×</button>
        </div>

        <div className="qr-modal-body">
          <div className="qr-code-container" ref={qrRef}>
            <QRCodeCanvas
              value={url}
              size={256}
              level="H"
              includeMargin={true}
            />
          </div>

          <div className="qr-info">
            <p className="qr-url">{url}</p>
            <p className="qr-short-code">Short Code: {shortCode}</p>
          </div>
        </div>

        <div className="qr-modal-actions">
          <button className="qr-action-button copy-button" onClick={copyToClipboard}>
            📋 Copy URL
          </button>
          <button className="qr-action-button download-button" onClick={downloadQR}>
            💾 Download QR
          </button>
        </div>
      </div>
    </div>
  );
};

export default QRCodeModal;