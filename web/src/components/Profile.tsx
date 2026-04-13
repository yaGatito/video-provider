import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import ProfileCard, { UserProfile } from './common/ProfileCard';
import VideoCard, { VideoPreview } from './common/VideoCard';
import Button from './common/Button';
import {
  PageShell,
  PageHeading,
  SectionBlock,
  SectionTitle,
  ContentGrid,
  Message,
} from './common/PageSection';

const PageCenter = styled(PageShell)`
  justify-items: center;
`;

const ErrorMessage = styled(Message)`
  width: min(100%, 460px);
`;

const VideosSection = styled(SectionBlock)`
  width: 100%;
  max-width: 980px;
`;

interface UserResponse {
  name: string;
  lastname: string;
  email: string;
  createdAt: string;
}

interface VideosResponse {
  videos: VideoPreview[];
}

const Profile: React.FC = () => {
  const [limit] = useState(10);
  const [offset] = useState(0);
  const [sort] = useState('date');
  const [order] = useState('t');
  const [user, setUser] = useState<UserProfile | null>(null);
  const [videos, setVideos] = useState<VideoPreview[]>([]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  const usersApiUrl = process.env.REACT_APP_USER_API_URL || '/userApi';
  const videoApiUrl = process.env.REACT_APP_VIDEO_API_URL || '/api';

  const decodeToken = (token: string): { user_id?: string; id?: string; sub?: string } | null => {
    try {
      const tokenPayload = token.split('.')[1];
      if (!tokenPayload) {
        return null;
      }
      const normalized = tokenPayload.replace(/-/g, '+').replace(/_/g, '/');
      const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=');
      const decodedJson = atob(padded);
      return JSON.parse(decodedJson);
    } catch {
      return null;
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    navigate('/login');
  };

  useEffect(() => {
    const fetchProfile = async () => {
      const token = localStorage.getItem('authToken');
      if (!token) {
        setError('Authentication token not found. Please log in.');
        setIsLoading(false);
        return;
      }

      const claims = decodeToken(token);
      const userId = claims?.user_id || claims?.id || claims?.sub;
      if (!userId) {
        setError('Unable to extract user ID from authentication token.');
        setIsLoading(false);
        return;
      }

      try {
        const response: AxiosResponse<UserResponse> = await axios.get(
          `${usersApiUrl}/v1/users/${userId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        setUser({
          ...response.data,
          id: userId,
          isCurrentUser: true,
        });

        const videosResponse: AxiosResponse<VideosResponse> = await axios.get(
          `${videoApiUrl}/v1/videos/pub/${userId}?&limit=${limit}&offset=${offset}&sort=${sort}&order=${order}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        setVideos(videosResponse.data.videos || []);
      } catch (fetchError: any) {
        setError(
          fetchError?.response?.data?.msg ||
            fetchError?.response?.data?.message ||
            'Unable to load profile information.',
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchProfile();
  }, [usersApiUrl, videoApiUrl, limit, order, sort]);

  if (isLoading) {
    return (
      <PageCenter>
        <PageHeading>User Profile</PageHeading>
        <p>Loading profile...</p>
      </PageCenter>
    );
  }

  if (error) {
    return (
      <PageCenter>
        <PageHeading>User Profile</PageHeading>
        <ErrorMessage $tone="error">{error}</ErrorMessage>
      </PageCenter>
    );
  }

  return (
    <PageCenter>
      <PageHeading>User Profile</PageHeading>
      {user ? (
        <>
          <ProfileCard
            user={user}
            showSubscribe={false}
          />
          <Button variant="secondary" onClick={handleLogout}>Logout</Button>
        </>
      ) : (
        <ErrorMessage $tone="error">No profile information found.</ErrorMessage>
      )}

      <VideosSection>
        <SectionTitle>Your uploads</SectionTitle>
        {videos.length > 0 ? (
          <ContentGrid>
            {videos.map((video) => (
              <VideoCard
                key={video.id}
                video={video}
                onClick={() => navigate(`/watch/${video.id}`)}
              />
            ))}
          </ContentGrid>
        ) : (
          <p>You have not uploaded any videos yet.</p>
        )}
      </VideosSection>
    </PageCenter>
  );
};

export default Profile;
