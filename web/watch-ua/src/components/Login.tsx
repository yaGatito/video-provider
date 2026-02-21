import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './Login.css';

interface UserLogin {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  message: string;
}

const Login: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    const apiUrl = process.env.REACT_APP_API_URL;
    axios.post<LoginResponse>(`${apiUrl}/api/login`, { username, password })
      .then((response: AxiosResponse<LoginResponse>) => {
        console.log('Login successful:', response.data);
        // Store token and redirect to home
        localStorage.setItem('authToken', response.data.token);
        window.location.href = '/';
      })
      .catch((error) => {
        console.error('There was an error logging in!', error);
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Invalid credentials. Please try again.'));
      });
  };

  return (
    <div className="login-page">
      <h1>🔐 Login</h1>
      {errorMessage && <p className="error-message">{errorMessage}</p>}
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Login</button>
      </form>
      <p className="signup-link">Don't have an account? <a href="/register">Sign up here</a></p>
    </div>
  );
};

export default Login;
