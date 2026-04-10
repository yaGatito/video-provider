import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';

interface UserLogin {
  email: string;
  password: string;
}

const Page = styled.div`
  display: grid;
  justify-items: center;
  gap: ${({ theme }) => theme.spacing.md};
`;

const Title = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const ErrorMessage = styled.p`
  width: min(100%, 460px);
  background: ${({ theme }) => theme.colors.errorBg};
  color: ${({ theme }) => theme.colors.errorText};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: ${({ theme }) => theme.spacing.md};
`;

const Form = styled.form`
  width: min(100%, 460px);
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  box-shadow: ${({ theme }) => theme.shadows.sm};
  padding: ${({ theme }) => theme.spacing.xl};
  display: grid;
  gap: ${({ theme }) => theme.spacing.md};
`;

const Input = styled.input`
  width: 100%;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.75rem 0.9rem;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

const Button = styled.button`
  border: none;
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.8rem 1rem;
  color: white;
  font-weight: 700;
  background: ${({ theme }) => theme.colors.brand};
  cursor: pointer;

  &:hover {
    background: ${({ theme }) => theme.colors.brandHover};
  }
`;

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
    <Page>
      <Title>Login</Title>
      {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
      <Form onSubmit={handleLogin}>
        <Input
          type="email"
          name="email"
          placeholder="Username"
          value={userData.email}
          onChange={handleChange}
          required
        />
        <Input
          type="password"
          name="password"
          placeholder="Password"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <Button type="submit">Login</Button>
      </Form>
    </Page>
  );
};

export default Login;
