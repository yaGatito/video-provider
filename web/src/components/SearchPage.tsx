import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';
import { useNavigate, useSearchParams } from 'react-router-dom';
import VideoCard, { VideoPreview } from './common/VideoCard';
import Button from './common/Button';
import {
  PageShell,
  PageHeading,
  FormShell,
  FormField,
  Label,
  TextInput,
  ContentGrid,
  Message,
} from './common/PageSection';

const Header = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const FilterRow = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.spacing.sm};
  align-items: center;
`;

const ToggleLabel = styled(Label)`
  display: inline-flex;
  gap: ${({ theme }) => theme.spacing.xs};
  padding: 0.65rem 0.85rem;
  background: ${({ theme }) => theme.colors.surface};
  cursor: pointer;
`;

const RadioInput = styled.input`
  accent-color: ${({ theme }) => theme.colors.brand};
`;

const Summary = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.95rem;
  margin-top: ${({ theme }) => theme.spacing.sm};
`;

const InfoMessage = styled(Message)`
  color: ${({ theme }) => theme.colors.textSecondary};
`;

interface VideosResp {
  videos: VideoPreview[];
}

const SearchPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const query = searchParams.get('query') ?? '';
  const scope = (searchParams.get('scope') === 'publishers' ? 'publishers' : 'videos') as 'videos' | 'publishers';
  const sortBy = (searchParams.get('sort') === 'publisher' ? 'publisher' : 'date') as 'date' | 'publisher';
  const [searchText, setSearchText] = useState(query);
  const [searchScope, setSearchScope] = useState<'videos' | 'publishers'>(scope);
  const [sortOrder, setSortOrder] = useState<'date' | 'publisher'>(sortBy);
  const [limit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [videos, setVideos] = useState<VideoPreview[]>([]);
  const videoApiUrl = process.env.REACT_APP_VIDEO_API_URL || '/api';
  const navigate = useNavigate();

  useEffect(() => {
    setSearchText(query);
    setSearchScope(scope);
    setSortOrder(sortBy);
  }, [query, scope, sortBy]);

  useEffect(() => {
    if (query.length > 2) {
      axios
        .get<VideosResp>(
          `${videoApiUrl}/v1/videos/search?query=${encodeURIComponent(query)}&scope=${scope}&sort=${sortBy}&order=t&limit=${limit}&offset=${offset}`,
        )
        .then((response: AxiosResponse<VideosResp>) => {
          setVideos(response.data.videos);
        })
        .catch((error: unknown) => {
          console.error('There was an error fetching the search results!', error);
        });
    } else {
      setVideos([]);
    }
  }, [query, scope, sortBy, limit, offset, videoApiUrl]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value);
  };

  const handleSearchSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const trimmed = searchText.trim();
    setOffset(0);
    setSearchParams({
      query: trimmed,
      scope: searchScope,
      sort: sortOrder,
    });
  };

  const loadMoreVideos = () => {
    setOffset((current) => current + limit);
  };

  return (
    <PageShell>
      <Header>
        <PageHeading>Search</PageHeading>
        <FormShell onSubmit={handleSearchSubmit}>
          <TextInput
            type="text"
            placeholder="Search for videos or publishers..."
            value={searchText}
            onChange={handleSearch}
          />
          <FilterRow>
            <ToggleLabel>
              <RadioInput
                type="radio"
                name="searchScope"
                value="videos"
                checked={searchScope === 'videos'}
                onChange={() => setSearchScope('videos')}
              />
              Videos
            </ToggleLabel>
            <ToggleLabel>
              <RadioInput
                type="radio"
                name="searchScope"
                value="publishers"
                checked={searchScope === 'publishers'}
                onChange={() => setSearchScope('publishers')}
              />
              Publishers
            </ToggleLabel>
            <ToggleLabel>
              <RadioInput
                type="radio"
                name="sortOrder"
                value="date"
                checked={sortOrder === 'date'}
                onChange={() => setSortOrder('date')}
              />
              Date
            </ToggleLabel>
            <ToggleLabel>
              <RadioInput
                type="radio"
                name="sortOrder"
                value="publisher"
                checked={sortOrder === 'publisher'}
                onChange={() => setSortOrder('publisher')}
              />
              Publisher
            </ToggleLabel>
          </FilterRow>
          <Button variant="secondary" type="submit">Search</Button>
        </FormShell>
        {query && (
          <Summary>
            Searching for <strong>{query}</strong> across{' '}
            {scope === 'publishers' ? 'publishers' : 'videos'} sorted by{' '}
            {sortBy === 'publisher' ? 'publisher' : 'date'}.
          </Summary>
        )}
      </Header>

      {videos.length > 0 ? (
        <>
          <ContentGrid>
            {videos.map((video) => (
              <VideoCard
                key={video.id}
                video={video}
                onClick={() => navigate(`/watch/${video.id}`)}
              />
            ))}
          </ContentGrid>
          <Button variant="secondary" onClick={loadMoreVideos}>Load More</Button>
        </>
      ) : (
        <InfoMessage>
          {query.length > 2 ? 'No videos found.' : 'Search for videos by typing at least 3 characters.'}
        </InfoMessage>
      )}
    </PageShell>
  );
};

export default SearchPage;
