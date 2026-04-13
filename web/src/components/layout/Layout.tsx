import React from 'react';
import styled from 'styled-components';
import Header from './Header';
import Footer from './Footer';
import FriendsSidebar from '../common/FriendsSidebar';

interface LayoutProps {
  children: React.ReactNode;
}

const Shell = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 100vh;
`;

const Main = styled.main`
  flex: 1;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: ${({ theme }) => theme.spacing.lg} ${({ theme }) => theme.spacing.sm};

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: ${({ theme }) => theme.spacing.md} ${({ theme }) => theme.spacing.sm};
  }
`;

const Content = styled.div`
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: ${({ theme }) => theme.spacing.lg};
  align-items: start;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
  }
`;

const PageContent = styled.div`
  min-width: 0;
`;

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <Shell>
      <Header />
      <Main>
        <Content>
          <FriendsSidebar />
          <PageContent>{children}</PageContent>
        </Content>
      </Main>
      <Footer />
    </Shell>
  );
};

export default Layout;
