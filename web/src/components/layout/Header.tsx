import React from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

const Wrapper = styled.header`
  position: sticky;
  top: 0;
  z-index: 1000;
  background: linear-gradient(90deg, #0c2842, #0f4c81);
  color: white;
  box-shadow: ${({ theme }) => theme.shadows.sm};
`;

const Container = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: ${({ theme }) => theme.spacing.md};
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: ${({ theme }) => theme.spacing.md};

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    flex-direction: column;
    align-items: flex-start;
  }
`;

const Logo = styled(Link)`
  font-family: ${({ theme }) => theme.fonts.heading};
  font-size: 1.35rem;
  font-weight: 700;
  letter-spacing: 0.04em;
`;

const Nav = styled.nav`
  width: 100%;
`;

const Menu = styled.ul`
  list-style: none;
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.spacing.lg};

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    justify-content: flex-start;
    gap: ${({ theme }) => theme.spacing.md};
  }
`;

const MenuLink = styled(Link)`
  color: #e6eef8;
  font-weight: 600;
  padding-bottom: 2px;
  border-bottom: 2px solid transparent;
  transition: color ${({ theme }) => theme.transitions.base}, border-color ${({ theme }) => theme.transitions.base};

  &:hover {
    color: #ffffff;
    border-bottom-color: ${({ theme }) => theme.colors.accent};
  }
`;

const Header: React.FC = () => {
  return (
    <Wrapper>
      <Container>
        <Logo to="/">Stream1</Logo>
        <Nav>
          <Menu>
            <li><MenuLink to="/">Home</MenuLink></li>
            <li><MenuLink to="/search">Search</MenuLink></li>
            <li><MenuLink to="/upload">Upload</MenuLink></li>
            <li><MenuLink to="/login">Login</MenuLink></li>
            <li><MenuLink to="/register">Register</MenuLink></li>
          </Menu>
        </Nav>
      </Container>
    </Wrapper>
  );
};

export default Header;
