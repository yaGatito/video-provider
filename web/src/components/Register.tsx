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

interface User {
  email: string;
  name: string;
  lastname: string;
  password: string;
}

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
    <CenteredPage>
      <PageHeading>Register</PageHeading>
      {errorMessage && <Message $tone="error">{errorMessage}</Message>}
      <FormShell onSubmit={handleRegister}>
        <TextInput
          type="text"
          name="name"
          placeholder="First Name"
          value={userData.name}
          onChange={handleChange}
          required
        />
        <TextInput
          type="text"
          name="lastname"
          placeholder="Last Name"
          value={userData.lastname}
          onChange={handleChange}
          required
        />
        <TextInput
          type="email"
          name="email"
          placeholder="Email"
          value={userData.email}
          onChange={handleChange}
          required
        />
        <TextInput
          type="password"
          name="password"
          placeholder="Password (min 8 characters)"
          value={userData.password}
          onChange={handleChange}
          required
        />
        <Button type="submit">Register</Button>
      </FormShell>
    </CenteredPage>
  );
};

export default Register;
