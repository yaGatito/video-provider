import React, { useState } from 'react';
import axios, { AxiosError, AxiosResponse } from 'axios';
import styled from 'styled-components';

interface VideoData {
  title: string;
  description?: string;
}

const Page = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.md};
  max-width: 700px;
  margin: 0 auto;
`;

const Title = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
`;

const Message = styled.p<{ $tone: 'success' | 'error' }>`
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: ${({ theme }) => theme.spacing.md};
  background: ${({ theme, $tone }) => ($tone === 'success' ? theme.colors.successBg : theme.colors.errorBg)};
  color: ${({ theme, $tone }) => ($tone === 'success' ? theme.colors.successText : theme.colors.errorText)};
`;

const Form = styled.form`
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  box-shadow: ${({ theme }) => theme.shadows.sm};
  padding: ${({ theme }) => theme.spacing.xl};
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
`;

const Group = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const Label = styled.label`
  font-weight: 700;
`;

const Input = styled.input`
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.75rem 0.9rem;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

const TextArea = styled.textarea`
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.75rem 0.9rem;
  min-height: 120px;
  resize: vertical;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

const SubmitButton = styled.button`
  border: none;
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.8rem 1rem;
  color: white;
  font-weight: 700;
  background: ${({ theme }) => theme.colors.brand};
  cursor: pointer;
  transition: background ${({ theme }) => theme.transitions.base};

  &:hover:not(:disabled) {
    background: ${({ theme }) => theme.colors.brandHover};
  }

  &:disabled {
    cursor: not-allowed;
    background: #8da8bf;
  }
`;

const UploadVideo: React.FC = () => {
  const [videoData, setVideoData] = useState<VideoData>({ title: '', description: '' });
  const [uploading, setUploading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const videoApiUrl = process.env.REACT_APP_VIDEO_API_URL || '/api';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setVideoData({ ...videoData, [name]: value });
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!videoData.title.trim()) {
      setErrorMessage('Please enter a video title');
      return;
    }

    if (!videoData.description || !videoData.description.trim()) {
      setErrorMessage('Please enter a video description');
      return;
    }

    if (videoData.title.length > 48) {
      setErrorMessage('Video title must not exceed 48 characters');
      return;
    }

    if (videoData.description.length > 512) {
      setErrorMessage('Video description must not exceed 512 characters');
      return;
    }

    setUploading(true);
    setErrorMessage('');
    setSuccessMessage('');

    const requestBody = {
      topic: videoData.title,
      description: videoData.description,
    };

    try {
      const response: AxiosResponse<{ message: string }> = await axios.post(
        `${videoApiUrl}/v1/videos/pub/123e4567-e89b-12d3-a456-426614174000`,
        requestBody
      );
      if (response.status >= 200 && response.status < 300) {
        setSuccessMessage(`Video "${videoData.title}" uploaded successfully!`);
      }
      setTimeout(() => {
        window.location.href = '/';
      }, 2000);
    } catch (error) {
      console.error('There was an error uploading the video!', error);
      let message = 'Please try again.';
      if (axios.isAxiosError(error)) {
        const axiosError: AxiosError<{ message: string }> = error;
        if (axiosError.response?.data?.message) {
          message = axiosError.response.data.message;
        }
      }
      setErrorMessage(`Error: ${message}`);
    } finally {
      setUploading(false);
    }
  };

  return (
    <Page>
      <Title>Upload Video</Title>
      {successMessage && <Message $tone="success">{successMessage}</Message>}
      {errorMessage && <Message $tone="error">{errorMessage}</Message>}

      <Form onSubmit={handleUpload}>
        <Group>
          <Label htmlFor="title">Video Title * (max 48 characters)</Label>
          <Input
            id="title"
            type="text"
            name="title"
            placeholder="Enter video title"
            value={videoData.title}
            onChange={handleChange}
            maxLength={48}
            required
          />
        </Group>

        <Group>
          <Label htmlFor="description">Description *</Label>
          <TextArea
            id="description"
            name="description"
            placeholder="Enter video description (max 512 characters)"
            value={videoData.description || ''}
            onChange={handleChange}
            rows={5}
            required
          />
        </Group>

        <SubmitButton type="submit" disabled={uploading}>
          {uploading ? 'Uploading...' : 'Upload Video'}
        </SubmitButton>
      </Form>
    </Page>
  );
};

export default UploadVideo;
