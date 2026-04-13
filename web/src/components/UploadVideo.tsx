import React, { useState } from 'react';
import axios, { AxiosError, AxiosResponse } from 'axios';
import Button from './common/Button';
import {
  CenteredPage,
  PageHeading,
  FormShell,
  FormField,
  Label,
  TextInput,
  TextAreaInput,
  Message,
} from './common/PageSection';

interface VideoData {
  title: string;
  description?: string;
}

const UploadVideo: React.FC = () => {
  const [videoData, setVideoData] = useState<VideoData>({ title: '', description: '' });
  const [uploading, setUploading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
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
      const token = localStorage.getItem('authToken');
      if (!token) {
        setErrorMessage('Authentication token not found. Please log in.');
        return;
      }

      const claims = decodeToken(token);
      const userId = claims?.user_id || claims?.id || claims?.sub;
      if (!userId) {
        setErrorMessage('Unable to extract user ID from authentication token.');
        return;
      }
      const response: AxiosResponse<{ message: string }> = await axios.post(
        `${videoApiUrl}/v1/videos/pub/${userId}`,
        requestBody,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
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
    <CenteredPage>
      <PageHeading>Upload Video</PageHeading>
      {successMessage && <Message $tone="success">{successMessage}</Message>}
      {errorMessage && <Message $tone="error">{errorMessage}</Message>}

      <FormShell onSubmit={handleUpload}>
        <FormField>
          <Label htmlFor="title">Video Title * (max 48 characters)</Label>
          <TextInput
            id="title"
            type="text"
            name="title"
            placeholder="Enter video title"
            value={videoData.title}
            onChange={handleChange}
            maxLength={48}
            required
          />
        </FormField>

        <FormField>
          <Label htmlFor="description">Description *</Label>
          <TextAreaInput
            id="description"
            name="description"
            placeholder="Enter video description (max 512 characters)"
            value={videoData.description || ''}
            onChange={handleChange}
            rows={5}
            required
          />
        </FormField>

        <Button type="submit" disabled={uploading}>
          {uploading ? 'Uploading...' : 'Upload Video'}
        </Button>
      </FormShell>
    </CenteredPage>
  );
};

export default UploadVideo;
