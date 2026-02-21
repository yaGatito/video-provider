import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './UploadVideo.css';

interface UploadResponse {
  id: number;
  title: string;
  message: string;
}

const UploadVideo: React.FC = () => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [previewImageFile, setPreviewImageFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleVideoFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setVideoFile(e.target.files[0]);
    }
  };

  const handlePreviewImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setPreviewImageFile(e.target.files[0]);
    }
  };

  const handleUpload = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) {
      setErrorMessage('Please enter a video title');
      return;
    }

    if (!videoFile) {
      setErrorMessage('Please select a video file');
      return;
    }

    setUploading(true);
    setErrorMessage('');
    setSuccessMessage('');

    const formData = new FormData();
    formData.append('title', title);
    formData.append('description', description);
    formData.append('video', videoFile);
    if (previewImageFile) {
      formData.append('previewImage', previewImageFile);
    }

    const apiUrl = process.env.REACT_APP_API_URL;
    axios.post<UploadResponse>(`${apiUrl}/api/videos/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
      .then((response: AxiosResponse<UploadResponse>) => {
        console.log('Video uploaded successfully:', response.data);
        setSuccessMessage(`Video "${response.data.title}" uploaded successfully!`);
        setTitle('');
        setDescription('');
        setVideoFile(null);
        setPreviewImageFile(null);
        // Reset form inputs
        const videoInput = document.getElementById('videoInput') as HTMLInputElement;
        const imageInput = document.getElementById('imageInput') as HTMLInputElement;
        if (videoInput) videoInput.value = '';
        if (imageInput) imageInput.value = '';
        // Redirect to home after 2 seconds
        setTimeout(() => {
          window.location.href = '/';
        }, 2000);
      })
      .catch((error) => {
        console.error('There was an error uploading the video!', error);
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Please try again.'));
      })
      .finally(() => {
        setUploading(false);
      });
  };

  return (
    <div className="upload-video-page">
      <h1>🎥 Upload Video</h1>
      
      {successMessage && <p className="success-message">{successMessage}</p>}
      {errorMessage && <p className="error-message">{errorMessage}</p>}
      
      <form onSubmit={handleUpload} className="upload-form">
        <div className="form-group">
          <label htmlFor="title">Video Title *</label>
          <input
            id="title"
            type="text"
            placeholder="Enter video title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            placeholder="Enter video description (optional)"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={5}
          />
        </div>

        <div className="form-group">
          <label htmlFor="videoInput">Video File *</label>
          <input
            id="videoInput"
            type="file"
            accept="video/*"
            onChange={handleVideoFileChange}
            required
          />
          {videoFile && <p className="file-name">Selected: {videoFile.name}</p>}
        </div>

        <div className="form-group">
          <label htmlFor="imageInput">Preview Image</label>
          <input
            id="imageInput"
            type="file"
            accept="image/*"
            onChange={handlePreviewImageChange}
          />
          {previewImageFile && <p className="file-name">Selected: {previewImageFile.name}</p>}
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
