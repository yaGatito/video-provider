import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './Login.css';

interface UserLogin {
  username: string;
  password: string;
}

const Login: React.FC = () => {
  const [userData, setUserData] = useState<UserLogin>({ username: '', password: '' });
  const [errorMessage, setErrorMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserData({ ...userData, [e.target.name]: e.target.value });
  };

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    axios.post('http://localhost:8081/v1/login', userData)
      .then((response: AxiosResponse<{ token: string; message: string }>) => {
        console.log('Login successful:', response.data);
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
          name="username"
          placeholder="Username"
          value={userData.username}
          onChange={handleChange}
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;