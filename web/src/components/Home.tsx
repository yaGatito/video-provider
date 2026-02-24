// src/components/Home.tsx

import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';

interface Video {
  id: number;
  title: string;
  description: string;
  likes: number;
  dislikes: number;
  comments: number;
  previewImage: string;
  uploadDate: string; // Add upload date
}

interface VideoResponse {
  videos: Video[];
}

const Home: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const apiUrl = process.env.VIDEO_SERVICE_API_URL;
    axios.get<VideoResponse>(`${apiUrl}/v1/videos/search?query=latest&limit=20&offset=0&sort=uploadDate&order=desc&`)
      .then((response: AxiosResponse<VideoResponse>) => {
        setVideos(response.data.videos);
        setLoading(false);
      })
      .catch((error: unknown) => {
        console.error('There was an error fetching the videos!', error);
        setError('Failed to load videos. Please try again later.');
        setLoading(false);
      });
  }, []);

  if (loading) {
    return (
      <div className="home">
        <div className="hero-section">
          <h1>Watch UA</h1>
          <p>Loading latest videos...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="home">
        <div className="hero-section">
          <h1>Watch UA</h1>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="home">
      <div className="hero-section">
        <h1>Watch UA</h1>
        <p>Discover amazing video content</p>
      </div>
      
      {/* Featured Video Section */}
      <div className="section">
        <h2>Featured Video</h2>
        <div className="video-card featured">
          <img 
            src="https://images.unsplash.com/photo-1618345522246-81e1f1e041b7?auto=format&fit=crop&w=1920&q=80" 
            alt="Featured Video" 
          />
          <div className="video-content">
            <h2>Sample Featured Video</h2>
            <p>This is a sample featured video to showcase the latest content on Watch UA.</p>
            <div className="video-stats">
              <span>👍 1.2k</span>
              <span>👎 120</span>
              <span>💬 450</span>
            </div>
          </div>
        </div>
      </div>
      
      {/* Latest Uploads Section */}
      <div className="section">
        <h2>Latest Uploads</h2>
        <div className="video-list">
          {videos.map(video => (
            <div key={video.id} className="video-card" onClick={() => window.location.href = `/v1/videos/id/${video.id}`}>
              <img src={video.previewImage} alt={video.title} />
              <div className="video-content">
                <h2>{video.title}</h2>
                <p>{video.description}</p>
                <div className="video-stats">
                  <span>👍 {video.likes}</span>
                  <span>👎 {video.dislikes}</span>
                  <span>💬 {video.comments}</span>
                  <span>📅 {video.uploadDate}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
      
      {/* Additional Content Section */}
      <div className="section">
        <h2>About Watch UA</h2>
        <p>
          Watch UA is a platform where you can discover amazing video content from creators around the world. 
          Whether you're looking for tutorials, entertainment, or educational content, we've got something for you.
        </p>
        <p>
          Join our community today and start exploring the world of video content!
        </p>
      </div>
    </div>
  );
};

export default Home;