import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';

interface User {
  email: string;
  name: string;
  lastname: string;
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
  width: min(100%, 480px);
  background: ${({ theme }) => theme.colors.errorBg};
  color: ${({ theme }) => theme.colors.errorText};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: ${({ theme }) => theme.spacing.md};
`;

const Form = styled.form`
  width: min(100%, 480px);
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

const Register: React.FC = () => {
  const [userData, setUserData] = useState<User>({ email: '', name: '', lastname: '', password: '' });
  const [errorMessage, setErrorMessage] = useState('');
  const userApiUrl = process.env.REACT_APP_USER_API_URL || '/userApi';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUserData({ ...userData, [e.target.name]: e.target.value });
  };

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    axios.post(`${userApiUrl}/v1/users`, userData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
      .then((response: AxiosResponse<User>) => {
        console.log('Registration successful:', response.data);
        window.location.href = '/home';
      })
      .catch((error) => {
        console.error('There was an error registering the user!', error);
        setErrorMessage(`Error: ${error.response?.data?.message || 'Please try again.'}`);
      });
  };

  return (
    <Page>
      <Title>Register</Title>
      {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
      <Form onSubmit={handleRegister}>
        <Input
          type="text"
          name="name"
          placeholder="First Name"
          value={userData.name}
          onChange={handleChange}
          required
        />
        <Input
          type="text"
          name="lastname"
          placeholder="Last Name"
          value={userData.lastname}
          onChange={handleChange}
          required
        />
        <Input
          type="email"
          name="email"
          placeholder="Email"
          value={userData.email}
          onChange={handleChange}
          required
        />
        <Input
          type="password"
          name="password"
          placeholder="Password (min 8 characters)"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <Button type="submit">Register</Button>
      </Form>
    </Page>
  );
};

export default Register;
