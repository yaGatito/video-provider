import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';

const Page = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
`;

const Header = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const Title = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const Input = styled.input`
  width: 100%;
  max-width: 680px;
  padding: 0.85rem 1rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  background: ${({ theme }) => theme.colors.surface};
  transition: border-color ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base};

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: ${({ theme }) => theme.spacing.lg};
`;

const Card = styled.article`
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  padding: ${({ theme }) => theme.spacing.lg};
  cursor: pointer;
  transition: transform ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base};

  &:hover {
    transform: translateY(-4px);
    box-shadow: ${({ theme }) => theme.shadows.md};
  }
`;

const Topic = styled.h2`
  font-family: ${({ theme }) => theme.fonts.heading};
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const Description = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: ${({ theme }) => theme.spacing.md};
`;

const Meta = styled.p`
  font-size: 0.85rem;
  color: ${({ theme }) => theme.colors.muted};
`;

const LoadMoreButton = styled.button`
  width: fit-content;
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: ${({ theme }) => theme.radius.sm};
  color: white;
  background: ${({ theme }) => theme.colors.brand};
  cursor: pointer;

  &:hover {
    background: ${({ theme }) => theme.colors.brandHover};
  }
`;

interface Video {
  id: string;
  publisherID: string;
  topic: string;
  description: string;
  createdAt: Date;
}

interface VideosResp {
  videos: Video[];
}

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [limit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [order] = useState('date');
  const [asc] = useState('t');
  const [videos, setVideos] = useState<Video[]>([]);
  const videoApiUrl = process.env.REACT_APP_VIDEO_API_URL || '/api';
  const navigate = useNavigate();

  useEffect(() => {
    if (query.length > 2) {
      axios.get<VideosResp>(`${videoApiUrl}/v1/videos/search?query=${query}&limit=${limit}&offset=${offset}&order=${order}&asc=${asc}`)
        .then((response: AxiosResponse<VideosResp>) => {
          setVideos(response.data.videos);
        })
        .catch((error: unknown) => {
          console.error('There was an error fetching the search results!', error);
        });
    }
  }, [query, limit, offset, order, asc, videoApiUrl]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
    setOffset(0);
  };

  const loadMoreVideos = () => {
    setOffset((current) => current + limit);
  };

  return (
    <Page>
      <Header>
        <Title>Search</Title>
        <Input
          type="text"
          placeholder="Search for videos..."
          value={query}
          onChange={handleSearch}
        />
      </Header>

      {videos.length > 0 ? (
        <>
          <Grid>
            {videos.map((video) => (
              <Card key={video.id} onClick={() => { navigate(`/watch/${video.id}`); }}>
                <Topic>{video.topic}</Topic>
                <Description>{video.description}</Description>
                <Meta>Published on {new Date(video.createdAt).toLocaleDateString()}</Meta>
              </Card>
            ))}
          </Grid>
          <LoadMoreButton onClick={loadMoreVideos}>Load More</LoadMoreButton>
        </>
      ) : (
        <Meta>No videos found.</Meta>
      )}
    </Page>
  );
};

export default SearchPage;
