import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './Register.css';

interface User {
  email: string;
  name: string;
  lastname: string;
  password: string;
}

const Register: React.FC = () => {
  const [userData, setUserData] = useState<User>({ email: '', name: '', lastname: '', password: '' });
  const [errorMessage, setErrorMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserData({ ...userData, [e.target.name]: e.target.value });
  };

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    axios.post('http://localhost:8081/v1/users', userData, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    })
      .then((response: AxiosResponse<User>) => {
        console.log('Registration successful:', response.data);
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
          name="name"
          placeholder="First Name"
          value={userData.name}
          onChange={handleChange}
          required
        />
        <input
          type="text"
          name="lastname"
          placeholder="Last Name"
          value={userData.lastname}
          onChange={handleChange}
          required
        />
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={userData.email}
          onChange={handleChange}
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password (min 8 characters)"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <button type="submit">Register</button>
      </form>
    </div>
  );
};

export default Register;