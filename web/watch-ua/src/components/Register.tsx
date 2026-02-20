import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './Register.css';

interface User {
  username: string;
  email: string;
  password: string;
}

const Register: React.FC = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    axios.post<User>('https://your-api-url.com/api/register', { username, email, password })
      .then((response: AxiosResponse<User>) => {
        console.log('Registration successful:', response.data);
        // Redirect to home or login page after registration
        window.location.href = '/home';
      })
      .catch((error) => {
        console.error('There was an error registering the user!', error);
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Please try again.'));
      });
  };

  return (
    <div className="register-page">
      <h1>📝 Register</h1>
      {errorMessage && <p className="error-message">{errorMessage}</p>}
      <form onSubmit={handleRegister}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Register</button>
      </form>
    </div>
  );
};

export default Register;