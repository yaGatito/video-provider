import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css';

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="header-container">
        <div className="logo">
          <Link to="/">🎥 Watch UA</Link>
        </div>
        <nav className="navbar">
          <ul className="nav-menu">
            <li><Link to="/">Home</Link></li>
            <li><Link to="/search">Search</Link></li>
            <li><Link to="/upload">Upload</Link></li>
            <li><Link to="/login">Login</Link></li>
            <li><Link to="/register">Register</Link></li>
          </ul>
        </nav>
      </div>
    </header>
  );
};

export default Header;
