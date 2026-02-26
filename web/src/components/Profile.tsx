import React from 'react';
import styled from 'styled-components';

interface User {
  username: string;
  email: string;
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

const Profile: React.FC<{ user: User }> = ({ user }) => {
  return (
    <Container>
      <Title>User Profile</Title>
      <Card>
        <Avatar src="/default-profile.png" alt="User Profile" />
        <Info>
          <Name>{user.username}</Name>
          <p>Email: {user.email}</p>
          <p>Member since: January 2023</p>
        </Info>
      </Card>
    </Container>
  );
};

export default Profile;
