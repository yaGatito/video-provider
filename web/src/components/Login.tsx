import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import Button from './common/Button';
import {
  CenteredPage,
  PageHeading,
  FormShell,
  TextInput,
  Message,
} from './common/PageSection';

interface UserLogin {
  email: string;
  password: string;
}

const Login: React.FC = () => {
  const [userData, setUserData] = useState<UserLogin>({ email: '', password: '' });
  const [errorMessage, setErrorMessage] = useState('');
  const usersApiUrl = process.env.REACT_APP_USER_API_URL || '/userApi';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserData({ ...userData, [e.target.name]: e.target.value });
  };

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    axios.post(`${usersApiUrl}/v1/login`, userData)
      .then((response: AxiosResponse<{ token: string; message: string }>) => {
        console.log('Login successful:', response.data);
        localStorage.setItem('authToken', response.data.token);
        window.location.href = '/';
      })
      .catch((error) => {
        console.error('There was an error logging in!', error);
        setErrorMessage(`Error: ${error.response?.data?.message || 'Invalid credentials. Please try again.'}`);
      });
  };

  return (
    <CenteredPage>
      <PageHeading>Login</PageHeading>
      {errorMessage && <Message $tone="error">{errorMessage}</Message>}
      <FormShell onSubmit={handleLogin}>
        <TextInput
          type="email"
          name="email"
          placeholder="Username"
          value={userData.email}
          onChange={handleChange}
          required
        />
        <TextInput
          type="password"
          name="password"
          placeholder="Password"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <Button type="submit">Login</Button>
      </FormShell>
    </CenteredPage>
  );
};

export default Login;
