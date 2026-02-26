import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';

interface Video {
  id: number;
  topic: string;
  description: string;
  previewImage: string;
  createdAt: string;
}

const Page = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.xl};
`;

const Hero = styled.section`
  background: linear-gradient(120deg, #0f4c81, #2a74ae);
  color: white;
  border-radius: ${({ theme }) => theme.radius.lg};
  padding: ${({ theme }) => theme.spacing.xxl} ${({ theme }) => theme.spacing.xl};
  box-shadow: ${({ theme }) => theme.shadows.md};
`;

const HeroTitle = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
  font-size: clamp(1.8rem, 3.5vw, 2.8rem);
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const HeroText = styled.p`
  max-width: 700px;
  color: #ddeeff;
`;

const Section = styled.section`
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
`;

const SectionTitle = styled.h2`
  font-family: ${({ theme }) => theme.fonts.heading};
  color: ${({ theme }) => theme.colors.textPrimary};
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: ${({ theme }) => theme.spacing.lg};
`;

const Card = styled.article<{ $featured?: boolean }>`
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  overflow: hidden;
  box-shadow: ${({ theme }) => theme.shadows.sm};
  transition: transform ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base};
  cursor: ${({ $featured }) => ($featured ? 'default' : 'pointer')};

  &:hover {
    transform: ${({ $featured }) => ($featured ? 'none' : 'translateY(-4px)')};
    box-shadow: ${({ theme }) => theme.shadows.md};
  }
`;

const Preview = styled.img`
  width: 100%;
  height: 210px;
  object-fit: cover;
  display: block;
`;

const CardBody = styled.div`
  padding: ${({ theme }) => theme.spacing.lg};
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const CardTitle = styled.h3`
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const CardText = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.95rem;
`;

const Meta = styled.div`
  font-size: 0.86rem;
  color: ${({ theme }) => theme.colors.muted};
`;

const Home: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const apiUrl = process.env.REACT_APP_API_URL;
  const navigate = useNavigate();

  const getRandomSearchQuery = (): string => {
    const queries = ['tec', 'wor', 'tra'];
    return queries[Math.floor(Math.random() * queries.length)];
  };

  useEffect(() => {
    axios.get<Video[]>(`${apiUrl}/v1/videos/search?limit=20&offset=0&orderBy=createdAt&asc=t&query=${getRandomSearchQuery()}`)
      .then((response: AxiosResponse<Video[]>) => {
        setVideos(response.data);
        setLoading(false);
      })
      .catch((requestError: unknown) => {
        console.error('There was an error fetching the videos!', requestError);
        setError('Failed to load videos. Please try again later.');
        setLoading(false);
      });
  }, [apiUrl]);

  if (loading) {
    return (
      <Page>
        <Hero>
          <HeroTitle>Watch UA</HeroTitle>
          <HeroText>Loading latest videos...</HeroText>
        </Hero>
      </Page>
    );
  }

  if (error) {
    return (
      <Page>
        <Hero>
          <HeroTitle>Watch UA</HeroTitle>
          <HeroText>{error}</HeroText>
        </Hero>
      </Page>
    );
  }

  return (
    <Page>
      <Hero>
        <HeroTitle>Watch UA</HeroTitle>
        <HeroText>Discover bold stories and fresh uploads from the community.</HeroText>
      </Hero>

      <Section>
        <SectionTitle>Featured Video</SectionTitle>
        <Card $featured>
          <Preview
            src="https://images.unsplash.com/photo-1618345522246-81e1f1e041b7?auto=format&fit=crop&w=1920&q=80"
            alt="Featured Video"
          />
          <CardBody>
            <CardTitle>Sample Featured Video</CardTitle>
            <CardText>This is a sample featured video to showcase the latest content on Watch UA.</CardText>
            <Meta>1.2k likes | 120 dislikes | 450 comments</Meta>
          </CardBody>
        </Card>
      </Section>

      <Section>
        <SectionTitle>Latest Uploads</SectionTitle>
        <Grid>
          {videos.map((video) => (
            <Card key={video.id} onClick={() => { navigate(`/watch/${video.id}`); }}>
              <Preview src={video.previewImage} alt={video.topic} />
              <CardBody>
                <CardTitle>{video.topic}</CardTitle>
                <CardText>{video.description}</CardText>
                <Meta>Published {new Date(video.createdAt).toLocaleDateString()}</Meta>
              </CardBody>
            </Card>
          ))}
        </Grid>
      </Section>
    </Page>
  );
};

export default Home;
