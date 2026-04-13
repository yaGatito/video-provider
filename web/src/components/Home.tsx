import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import VideoCard, { VideoPreview } from './common/VideoCard';
import { CenteredPage, SectionBlock, SectionTitle, ContentGrid } from './common/PageSection';

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

interface VideosResp {
  videos: VideoPreview[];
}

const Home: React.FC = () => {
  const [videos, setVideos] = useState<VideoPreview[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const apiUrl = process.env.REACT_APP_VIDEO_API_URL || '/api';
  const navigate = useNavigate();

  const getRandomSearchQuery = (): string => {
    const queries = ['tec', 'wor', 'tra'];
    return queries[Math.floor(Math.random() * queries.length)];
  };

  useEffect(() => {
    axios
      .get<VideosResp>(
        `${apiUrl}/v1/videos/search?limit=20&offset=0&sort=date&order=t&query=${getRandomSearchQuery()}`,
      )
      .then((response: AxiosResponse<VideosResp>) => {
        setVideos(response.data.videos);
      })
      .catch((requestError: unknown) => {
        console.error('There was an error fetching the videos!', requestError);
        setError('Failed to load videos. Please try again later.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [apiUrl]);

  if (loading) {
    return (
      <CenteredPage>
        <Hero>
          <HeroTitle>Stream</HeroTitle>
          <HeroText>Loading latest videos...</HeroText>
        </Hero>
      </CenteredPage>
    );
  }

  if (error) {
    return (
      <CenteredPage>
        <Hero>
          <HeroTitle>Stream1</HeroTitle>
          <HeroText>{error}</HeroText>
        </Hero>
      </CenteredPage>
    );
  }

  return (
    <CenteredPage>
      <Hero>
        <HeroTitle>Stream1</HeroTitle>
        <HeroText>Discover bold stories and fresh uploads from the community.</HeroText>
      </Hero>

      <SectionBlock>
        <SectionTitle>Featured Video</SectionTitle>
        <VideoCard
          video={{
            id: 'featured',
            topic: 'Sample Featured Video',
            description: 'This is a sample featured video to showcase the latest content on Stream1.',
            createdAt: new Date().toISOString(),
          }}
          onClick={() => navigate('/watch/featured')}
        />
      </SectionBlock>

      <SectionBlock>
        <SectionTitle>Latest Uploads</SectionTitle>
        <ContentGrid>
          {videos.map((video) => (
            <VideoCard
              key={video.id}
              video={video}
              onClick={() => navigate(`/watch/${video.id}`)}
            />
          ))}
        </ContentGrid>
      </SectionBlock>
    </CenteredPage>
  );
};

export default Home;
