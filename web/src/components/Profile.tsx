import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';

interface User {
  name: string;
  lastname: string;
  email: string;
  createdAt: string;
}

const Container = styled.section`
  display: grid;
  justify-items: center;
  gap: ${({ theme }) => theme.spacing.lg};
`;

const Title = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const Card = styled.article`
  width: min(100%, 460px);
  display: grid;
  justify-items: center;
  gap: ${({ theme }) => theme.spacing.md};
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  box-shadow: ${({ theme }) => theme.shadows.sm};
  padding: ${({ theme }) => theme.spacing.xl};
`;

const Avatar = styled.img`
  width: 130px;
  height: 130px;
  border-radius: 50%;
  object-fit: cover;
`;

const Info = styled.div`
  width: 100%;
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const Name = styled.h2`
  color: ${({ theme }) => theme.colors.textPrimary};
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const ErrorMessage = styled.p`
  width: min(100%, 460px);
  background: ${({ theme }) => theme.colors.errorBg};
  color: ${({ theme }) => theme.colors.errorText};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: ${({ theme }) => theme.spacing.md};
`;

const Profile: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const usersApiUrl = process.env.REACT_APP_USER_API_URL || '/userApi';

  const decodeToken = (token: string): { user_id?: string; id?: string; sub?: string } | null => {
    try {
      const tokenPayload = token.split('.')[1];
      if (!tokenPayload) {
        return null;
      }
      const normalized = tokenPayload.replace(/-/g, '+').replace(/_/g, '/');
      const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=');
      const decodedJson = atob(padded);
      return JSON.parse(decodedJson);
    } catch {
      return null;
    }
  };

  useEffect(() => {
    const fetchProfile = async () => {
      const token = localStorage.getItem('authToken');
      if (!token) {
        setError('Authentication token not found. Please log in.');
        setIsLoading(false);
        return;
      }

      const claims = decodeToken(token);
      const userId = claims?.user_id || claims?.id || claims?.sub;
      if (!userId) {
        setError('Unable to extract user ID from authentication token.');
        setIsLoading(false);
        return;
      }

      try {
        const response: AxiosResponse<User> = await axios.get(
          `${usersApiUrl}/v1/users/${userId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        setUser(response.data);
      } catch (fetchError: any) {
        setError(
          fetchError.response?.data?.msg ||
            fetchError.response?.data?.message ||
            'Unable to load profile information.',
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchProfile();
  }, [usersApiUrl]);

  if (isLoading) {
    return (
      <Container>
        <Title>User Profile</Title>
        <p>Loading profile...</p>
      </Container>
    );
  }

  if (error) {
    return (
      <Container>
        <Title>User Profile</Title>
        <ErrorMessage>{error}</ErrorMessage>
      </Container>
    );
  }

  if (!user) {
    return (
      <Container>
        <Title>User Profile</Title>
        <p>No profile information available.</p>
      </Container>
    );
  }

  return (
    <Container>
      <Title>User Profile</Title>
      <Card>
        <Avatar src="/default-profile.png" alt="User Profile" />
        <Info>
          <Name>{`${user.name} ${user.lastname}`}</Name>
          <p>Email: {user.email}</p>
          <p>Member since: {new Date(user.createdAt).toLocaleDateString()}</p>
        </Info>
      </Card>
    </Container>
  );
};

export default Profile;
