import React from 'react';
import styled from 'styled-components';

const Wrapper = styled.footer`
  margin-top: ${({ theme }) => theme.spacing.xxl};
  background: #0c2842;
  color: #d2dfef;
  border-top: 1px solid #244566;
`;

const Container = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: ${({ theme }) => theme.spacing.xl} ${({ theme }) => theme.spacing.md};
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: ${({ theme }) => theme.spacing.xl};
`;

const Title = styled.h4`
  margin-bottom: ${({ theme }) => theme.spacing.md};
  font-family: ${({ theme }) => theme.fonts.heading};
  color: #ffffff;
`;

const Text = styled.p`
  color: #c7d6e9;
  font-size: 0.95rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const List = styled.ul`
  list-style: none;
`;

const ListItem = styled.li`
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const FooterLink = styled.a`
  color: #d2dfef;
  transition: color ${({ theme }) => theme.transitions.base};

  &:hover {
    color: #ffffff;
  }
`;

const Socials = styled.div`
  margin-top: ${({ theme }) => theme.spacing.md};
  display: flex;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const SocialIcon = styled.a`
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #1a456f;
  transition: transform ${({ theme }) => theme.transitions.base}, background ${({ theme }) => theme.transitions.base};

  &:hover {
    background: ${({ theme }) => theme.colors.accent};
    transform: translateY(-2px);
  }
`;

const Bottom = styled.div`
  margin-top: ${({ theme }) => theme.spacing.xl};
  padding-top: ${({ theme }) => theme.spacing.md};
  border-top: 1px solid #244566;
  text-align: center;
  color: #9eb4cc;
  font-size: 0.85rem;
`;

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <Wrapper>
      <Container>
        <Grid>
          <section>
            <Title>About</Title>
            <Text>Stream1 is a platform for discovering and publishing Ukrainian video content.</Text>
          </section>

          <section>
            <Title>Quick Links</Title>
            <List>
              <ListItem><FooterLink href="/">Home</FooterLink></ListItem>
              <ListItem><FooterLink href="/search">Search</FooterLink></ListItem>
              <ListItem><FooterLink href="/upload">Upload</FooterLink></ListItem>
            </List>
          </section>

          <section>
            <Title>Contact</Title>
            <Text>Email: <FooterLink href="mailto:info@stream1.com">info@stream1.com</FooterLink></Text>
            <Text>Support: <FooterLink href="mailto:support@stream1.com">support@stream1.com</FooterLink></Text>
            <Text>Phone: <FooterLink href="tel:+380441234567">+38 (044) 123-45-67</FooterLink></Text>
            <Text>Address: Odesa, Ukraine</Text>
            <Socials>
              <SocialIcon href="https://facebook.com/stream1" title="Facebook" aria-label="Facebook">f</SocialIcon>
              <SocialIcon href="https://twitter.com/stream1" title="Twitter" aria-label="Twitter">X</SocialIcon>
              <SocialIcon href="https://instagram.com/stream1" title="Instagram" aria-label="Instagram">I</SocialIcon>
              <SocialIcon href="https://youtube.com/@stream1" title="YouTube" aria-label="YouTube">Y</SocialIcon>
            </Socials>
          </section>
        </Grid>

        <Bottom>&copy; {currentYear} Stream1. All rights reserved.</Bottom>
      </Container>
    </Wrapper>
  );
};

export default Footer;
