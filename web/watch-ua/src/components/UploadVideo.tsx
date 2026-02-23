import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './UploadVideo.css';

interface VideoData {
  title: string;
  description?: string;
}

const UploadVideo: React.FC = () => {
  const [videoData, setVideoData] = useState<VideoData>({ title: '', description: '' });
  const [uploading, setUploading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

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
      await axios.post('http://localhost:8080/v1/videos/pub/123e4567-e89b-12d3-a456-426614174000', requestBody);
      setSuccessMessage(`Video "${videoData.title}" uploaded successfully!`);
      setTimeout(() => {
        window.location.href = '/';
      }, 2000);
    } catch (error) {
      console.error('There was an error uploading the video!', error);
      const message = axios.isAxiosError(error) && error.response?.data?.message ? error.response.data.message : 'Please try again.';
      setErrorMessage('Error: ' + message);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="upload-video-page">
      <h1>🎥 Upload Video</h1>

      {successMessage && <p className="success-message">{successMessage}</p>}
      {errorMessage && <p className="error-message">{errorMessage}</p>}

      <form onSubmit={handleUpload} className="upload-form">
        <div className="form-group">
          <label htmlFor="title">Video Title * (max 48 characters)</label>
          <input
            id="title"
            type="text"
            name="title"
            placeholder="Enter video title"
            value={videoData.title}
            onChange={handleChange}
            maxLength={48}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description *</label>
          <textarea
            id="description"
            name="description"
            placeholder="Enter video description (max 512 characters)"
            value={videoData.description || ''}
            onChange={handleChange}
            rows={5}
            required
          />
        </div>

        <button
          type="submit"
          className="upload-button"
          disabled={uploading}
        >
          {uploading ? 'Uploading...' : 'Upload Video'}
        </button>
      </form>
    </div>
  );
};

export default UploadVideo;
