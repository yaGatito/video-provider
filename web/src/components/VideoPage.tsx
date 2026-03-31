import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import { Link, useParams } from 'react-router-dom';
import styled from 'styled-components';

interface Video {
  id: string;
  title?: string;
  topic?: string;
  description: string;
  likes?: number;
  dislikes?: number;
  comments?: number;
  previewImage?: string;
  createdAt?: string;
}

const Page = styled.div`
  display: grid;
  justify-items: center;
`;

const Container = styled.article`
  width: min(100%, 900px);
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  overflow: hidden;
  box-shadow: ${({ theme }) => theme.shadows.md};
`;

const Header = styled.header`
  display: grid;
  gap: ${({ theme }) => theme.spacing.md};
  background: ${({ theme }) => theme.colors.surfaceAlt};
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  padding: ${({ theme }) => theme.spacing.xl};
`;

const BackButton = styled(Link)`
  width: fit-content;
  background: ${({ theme }) => theme.colors.brand};
  color: white;
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.55rem 0.85rem;
  font-weight: 600;

  &:hover {
    background: ${({ theme }) => theme.colors.brandHover};
  }
`;

const Title = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
  font-size: clamp(1.5rem, 3.2vw, 2.2rem);
`;

const Body = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
  padding: ${({ theme }) => theme.spacing.xl};
`;

const Image = styled.img`
  width: 100%;
  border-radius: ${({ theme }) => theme.radius.sm};
`;

const VideoPlayer = styled.video`
  width: 100%;
  border-radius: ${({ theme }) => theme.radius.sm};
  background: #000;
`;

const Description = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  white-space: pre-line;
`;

const Stats = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.spacing.md};
`;

const Stat = styled.span`
  background: ${({ theme }) => theme.colors.surfaceAlt};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.45rem 0.7rem;
  font-size: 0.92rem;
`;

const StreamHint = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.92rem;
`;

const Loading = styled.div`
  font-size: 1.1rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const VideoPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [video, setVideo] = useState<Video | null>(null);
  const [streamUnavailable, setStreamUnavailable] = useState(false);

  useEffect(() => {
    const apiUrl = process.env.REACT_APP_API_URL;
    axios.get<Video>(`${apiUrl}/v1/videos/id/${id}`)
      .then((response: AxiosResponse<Video>) => {
        setVideo(response.data);
      })
      .catch((error: unknown) => {
        console.error('There was an error fetching the video!', error);
      });
  }, [id]);

  if (!video) {
    return <Loading>Loading video...</Loading>;
  }

  const apiUrl = process.env.REACT_APP_API_URL || '';
  const streamBase = process.env.REACT_APP_VIDEO_STREAM_URL || `${apiUrl}/v1/videos/stream`;
  const streamUrl = `${streamBase}/${id}`;
  const title = video.topic || video.title || `Video ${video.id}`;

  return (
    <Page>
      <Container>
        <Header>
          <BackButton to="/">Back to Videos</BackButton>
          <Title>{title}</Title>
        </Header>
        <Body>
          {!streamUnavailable ? (
            <VideoPlayer
              controls
              preload="metadata"
              poster={video.previewImage}
              onError={() => setStreamUnavailable(true)}
            >
              <source src={streamUrl} />
            </VideoPlayer>
          ) : (
            <Image src={video.previewImage || 'https://placehold.co/1200x675?text=Video+Unavailable'} alt={title} />
          )}
          {streamUnavailable && (
            <StreamHint>
              Stream1 is not available from <code>{streamUrl}</code>. Set <code>REACT_APP_VIDEO_STREAM_URL</code> if your
              stream endpoint is different.
            </StreamHint>
          )}
          <Description>{video.description}</Description>
          <Stats>
            <Stat>{video.likes ?? 0} Likes</Stat>
            <Stat>{video.dislikes ?? 0} Dislikes</Stat>
            <Stat>{video.comments ?? 0} Comments</Stat>
            {video.createdAt && <Stat>{new Date(video.createdAt).toLocaleDateString()}</Stat>}
          </Stats>
        </Body>
      </Container>
    </Page>
  );
};

export default VideoPage;
