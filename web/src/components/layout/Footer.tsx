import React from 'react';
import './Footer.css';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="footer">
      <div className="footer-container">
        <div className="footer-content">
          <div className="footer-section">
            <h4>About</h4>
            <p>Watch UA - Your premier platform for Ukrainian video content.</p>
          </div>
          <div className="footer-section">
            <h4>Quick Links</h4>
            <ul>
              <li><a href="/">Home</a></li>
              <li><a href="/search">Search</a></li>
              <li><a href="/upload">Upload</a></li>
            </ul>
          </div>
          <div className="footer-section">
            <h4>Contact</h4>
            <p>Email: <a href="mailto:info@watchua.com">info@watchua.com</a></p>
            <p>Support: <a href="mailto:support@watchua.com">support@watchua.com</a></p>
            <p>Phone: <a href="tel:+380441234567">+38 (044) 123-45-67</a></p>
            <p>Address: Kyiv, Ukraine</p>
            <div className="social-links">
              <a href="https://facebook.com/watchua" className="social-icon" title="Facebook">f</a>
              <a href="https://twitter.com/watchua" className="social-icon" title="Twitter">𝕏</a>
              <a href="https://instagram.com/watchua" className="social-icon" title="Instagram">📷</a>
              <a href="https://youtube.com/@watchua" className="social-icon" title="YouTube">▶</a>
            </div>
          </div>
        </div>
        <div className="footer-bottom">
          <p>&copy; {currentYear} Watch UA. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
