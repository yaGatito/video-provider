import React, { useEffect, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styled, { css } from 'styled-components';
import Button from '../common/Button';

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

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    flex-direction: column;
    align-items: stretch;
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

const SearchWrapper = styled.div`
  position: relative;
  display: flex;
  align-items: center;
  width: min(100%, 420px);

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    width: 100%;
  }
`;

const SearchField = styled.input`
  flex: 1;
  min-width: 160px;
  padding: 0.75rem 1rem;
  padding-right: 108px;
  border-radius: ${({ theme }) => theme.radius.lg};
  border: 1px solid rgba(15, 76, 129, 0.12);
  background: ${({ theme }) => theme.colors.surface};
  background-image: linear-gradient(180deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.8));
  color: ${({ theme }) => theme.colors.textPrimary};
  transition: border-color ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base}, background ${({ theme }) => theme.transitions.base};

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 4px rgba(15, 76, 129, 0.1);
    background-image: linear-gradient(180deg, rgba(255, 255, 255, 1), rgba(248, 250, 255, 0.9));
  }
`;

const SearchButton = styled(Button).attrs({
  variant: 'secondary',
  type: 'submit',
})`
  position: absolute;
  right: 4px;
  top: 4px;
  bottom: 4px;
  width: 90px;
  border-radius: 0 ${({ theme }) => theme.radius.lg} ${({ theme }) => theme.radius.lg} 0;
  padding: 0 0.9rem;
  min-height: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background-image: linear-gradient(135deg, rgba(15, 76, 129, 0.95), rgba(11, 58, 99, 0.95));
  color: white;
  box-shadow: none;
  transition: background ${({ theme }) => theme.transitions.base}, opacity ${({ theme }) => theme.transitions.base};

  &:hover {
    background-image: linear-gradient(135deg, rgba(15, 76, 129, 1), rgba(11, 58, 99, 1));
    opacity: 0.96;
  }
`;

const FilterPanel = styled.div<{ $visible: boolean }>`
  position: absolute;
  top: calc(100% + 0.65rem);
  left: 0;
  z-index: 10;
  width: 100%;
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  box-shadow: ${({ theme }) => theme.shadows.md};
  overflow: hidden;
  opacity: ${({ $visible }) => ($visible ? 1 : 0)};
  transform: translateY(${({ $visible }) => ($visible ? '0' : '-10px')});
  max-height: ${({ $visible }) => ($visible ? '240px' : '0')};
  padding: ${({ $visible, theme }) => ($visible ? theme.spacing.md : '0')};
  transition: opacity 220ms ease, transform 220ms ease, max-height 220ms ease, padding 220ms ease;
`;

const FilterGroup = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const FilterTitle = styled.h4`
  margin: 0;
  font-size: 0.95rem;
  color: ${({ theme }) => theme.colors.textPrimary};
`;

const FilterOption = styled.label`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.sm};
  padding: 0.65rem;
  border-radius: ${({ theme }) => theme.radius.sm};
  transition: background ${({ theme }) => theme.transitions.base};

  &:hover {
    background: ${({ theme }) => theme.colors.surfaceAlt};
  }
`;

const FilterInput = styled.input`
  accent-color: ${({ theme }) => theme.colors.brand};
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

const navItemStyles = css`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
  padding: 0.75rem 1rem;
  color: #e6eef8;
  font-weight: 600;
  position: relative;
  background: transparent;
  text-decoration: none;
  transition: color ${({ theme }) => theme.transitions.base}, background ${({ theme }) => theme.transitions.base};

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 2px;
    background: ${({ theme }) => theme.colors.accent};
    transform: scaleX(0);
    transform-origin: left center;
    transition: transform ${({ theme }) => theme.transitions.base};
  }

  &:hover {
    color: #ffffff;
  }

  &:hover::after {
    transform: scaleX(1);
  }

  &:focus-visible {
    outline: 2px solid rgba(255, 255, 255, 0.35);
    outline-offset: 2px;
  }
`;

const MenuLink = styled(Link)`
  ${navItemStyles}
`;

const MenuButton = styled.button`
  ${navItemStyles}
  appearance: none;
  border: none;
  box-shadow: none;
  background: transparent;
  cursor: pointer;
  font: inherit;
`;

const Header: React.FC = () => {
  const [isAuthorized, setIsAuthorized] = useState<boolean>(Boolean(localStorage.getItem('authToken')));
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [searchScope, setSearchScope] = useState<'videos' | 'publishers'>('videos');
  const [sortBy, setSortBy] = useState<'date' | 'publisher'>('date');
  const [showFilters, setShowFilters] = useState<boolean>(false);
  const searchRef = useRef<HTMLDivElement | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const handleStorage = () => {
      setIsAuthorized(Boolean(localStorage.getItem('authToken')));
    };

    window.addEventListener('storage', handleStorage);
    return () => window.removeEventListener('storage', handleStorage);
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (!searchRef.current?.contains(event.target as Node)) {
        setShowFilters(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    setIsAuthorized(false);
    navigate('/login');
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmedQuery = searchQuery.trim();
    if (!trimmedQuery) {
      return;
    }

    setShowFilters(false);
    navigate(
      `/search?query=${encodeURIComponent(trimmedQuery)}&scope=${searchScope}&sort=${sortBy}`,
    );
  };

  return (
    <Wrapper>
      <Container>
        <Logo to="/">Stream1</Logo>
        <SearchWrapper ref={searchRef}>
          <form onSubmit={handleSubmit} style={{ position: 'relative', width: '100%' }}>
            <SearchField
              aria-label="Search videos or publishers"
              placeholder="Search videos or publishers"
              value={searchQuery}
              onChange={(event) => setSearchQuery(event.target.value)}
              onFocus={() => setShowFilters(true)}
            />
            <SearchButton>Search</SearchButton>
          </form>
          <FilterPanel $visible={showFilters}>
            <FilterGroup>
              <FilterTitle>Search in</FilterTitle>
              <FilterOption>
                <FilterInput
                  type="radio"
                  name="searchScope"
                  value="videos"
                  checked={searchScope === 'videos'}
                  onChange={() => setSearchScope('videos')}
                />
                Videos
              </FilterOption>
              <FilterOption>
                <FilterInput
                  type="radio"
                  name="searchScope"
                  value="publishers"
                  checked={searchScope === 'publishers'}
                  onChange={() => setSearchScope('publishers')}
                />
                Publishers
              </FilterOption>
            </FilterGroup>
            <FilterGroup>
              <FilterTitle>Sort by</FilterTitle>
              <FilterOption>
                <FilterInput
                  type="radio"
                  name="sortBy"
                  value="date"
                  checked={sortBy === 'date'}
                  onChange={() => setSortBy('date')}
                />
                Date
              </FilterOption>
              <FilterOption>
                <FilterInput
                  type="radio"
                  name="sortBy"
                  value="publisher"
                  checked={sortBy === 'publisher'}
                  onChange={() => setSortBy('publisher')}
                />
                Publisher
              </FilterOption>
            </FilterGroup>
          </FilterPanel>
        </SearchWrapper>
        <Nav>
          <Menu>
            <li><MenuLink to="/">Home</MenuLink></li>
            <li><MenuLink to="/upload">Upload</MenuLink></li>
            <li><MenuLink to="/profile">Profile</MenuLink></li>
            {!isAuthorized && <li><MenuLink to="/login">Login</MenuLink></li>}
            {!isAuthorized && <li><MenuLink to="/register">Register</MenuLink></li>}
            {isAuthorized && <li><MenuButton type="button" onClick={handleLogout}>Logout</MenuButton></li>}
          </Menu>
        </Nav>
      </Container>
    </Wrapper>
  );
};

export default Header;
