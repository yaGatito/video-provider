import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';
import Button from './Button';

interface Friend {
  id: string;
  name: string;
  watching: string;
  isOnline: boolean;
}

const Sidebar = styled.aside`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
  width: 100%;
  max-width: 240px;
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  padding: ${({ theme }) => theme.spacing.md};
  box-shadow: ${({ theme }) => theme.shadows.sm};
  position: sticky;
  top: calc(${({ theme }) => theme.spacing.lg} + 56px);
  align-self: start;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    position: static;
    max-width: 100%;
    top: 0;
  }
`;

const SidebarTitle = styled.h2`
  font-family: ${({ theme }) => theme.fonts.heading};
  font-size: 1rem;
  margin: 0;
  color: ${({ theme }) => theme.colors.textPrimary};
`;

const FriendCard = styled.article`
  display: grid;
  gap: ${({ theme }) => theme.spacing.xs};
  padding: ${({ theme }) => theme.spacing.sm};
  border-radius: ${({ theme }) => theme.radius.sm};
  background: ${({ theme }) => theme.colors.surfaceAlt};
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const FriendName = styled.h3`
  margin: 0;
  font-size: 0.95rem;
  color: ${({ theme }) => theme.colors.textPrimary};
`;

const FriendWatching = styled.p`
  margin: 0;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  line-height: 1.35;
`;

const Status = styled.span`
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.85rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const StatusDot = styled.span<{ $online: boolean }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: ${({ theme, $online }) => ($online ? '#22c55e' : theme.colors.muted)};
`;

const EmptyState = styled.p`
  margin: 0;
  color: ${({ theme }) => theme.colors.textSecondary};
  line-height: 1.6;
`;

const FooterNote = styled.p`
  margin: 0;
  font-size: 0.9rem;
  color: ${({ theme }) => theme.colors.muted};
`;

const mockFriends: Friend[] = [
  { id: 'f1', name: 'Ariana', watching: 'Design thinking for creators', isOnline: true },
  { id: 'f2', name: 'Mauro', watching: 'Stream1 release notes', isOnline: true },
  { id: 'f3', name: 'Leila', watching: 'React hooks deep dive', isOnline: false },
  { id: 'f4', name: 'Noah', watching: 'Fast upload workflows', isOnline: true },
];

const FriendsSidebar: React.FC = () => {
  const [isAuthorized, setIsAuthorized] = useState<boolean>(Boolean(localStorage.getItem('authToken')));

  useEffect(() => {
    const handleStorage = () => {
      setIsAuthorized(Boolean(localStorage.getItem('authToken')));
    };

    window.addEventListener('storage', handleStorage);
    return () => window.removeEventListener('storage', handleStorage);
  }, []);

  return (
    <Sidebar>
      <SidebarTitle>Friends Watching Now</SidebarTitle>
      {isAuthorized ? (
        mockFriends.map((friend) => (
          <FriendCard key={friend.id}>
            <FriendName>{friend.name}</FriendName>
            <FriendWatching>Watching <strong>{friend.watching}</strong></FriendWatching>
            <Status>
              <StatusDot $online={friend.isOnline} />
              {friend.isOnline ? 'Online' : 'Offline'}
            </Status>
          </FriendCard>
        ))
      ) : (
        <EmptyState>Login to see what your friends are watching and explore content together.</EmptyState>
      )}
      <FooterNote>Friends list is available across every page once you sign in.</FooterNote>
    </Sidebar>
  );
};

export default FriendsSidebar;
