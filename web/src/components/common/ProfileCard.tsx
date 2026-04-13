import React from 'react';
import styled from 'styled-components';
import Button from './Button';

export interface UserProfile {
  name: string;
  lastname: string;
  email: string;
  createdAt: string;
  id?: string;
  isCurrentUser?: boolean;
}

interface ProfileCardProps {
  user: UserProfile;
  onSubscribe?: () => void;
  showSubscribe?: boolean;
}

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
  text-align: center;
`;

const Name = styled.h2`
  color: ${({ theme }) => theme.colors.textPrimary};
  font-family: ${({ theme }) => theme.fonts.heading};
  margin: 0;
`;

const Label = styled.span`
  display: block;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.95rem;
`;

const Badge = styled.div`
  width: 100%;
  padding: ${({ theme }) => theme.spacing.sm};
  border-radius: ${({ theme }) => theme.radius.sm};
  background: ${({ theme }) => theme.colors.surfaceAlt};
  border: 1px solid ${({ theme }) => theme.colors.border};
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.88rem;
`;

const ActionRow = styled.div`
  display: grid;
  width: 100%;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const ProfileCard: React.FC<ProfileCardProps> = ({ user, onSubscribe, showSubscribe }) => {
  return (
    <Card>
      <Avatar src="/default-profile.png" alt="User profile" />
      <Info>
        <Name>{`${user.name} ${user.lastname}`}</Name>
        <Label>Email: {user.email}</Label>
        <Label>Member since: {new Date(user.createdAt).toLocaleDateString()}</Label>
      </Info>

      {user.isCurrentUser && user.id && (
        <Badge>Your user ID: {user.id}</Badge>
      )}

      {showSubscribe && onSubscribe && !user.isCurrentUser && (
        <ActionRow>
          <Button variant="secondary" onClick={onSubscribe}>
            Subscribe
          </Button>
        </ActionRow>
      )}
    </Card>
  );
};

export default ProfileCard;
