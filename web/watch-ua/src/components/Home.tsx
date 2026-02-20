import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './Home.css';

interface Video {
  id: number;
  title: string;
  description: string;
  likes: number;
  dislikes: number;
  comments: number;
  previewImage: string;
}

interface VideoResponse {
  videos: Video[];
}

const Home: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    const apiUrl = process.env.REACT_APP_API_URL;
    axios.get<VideoResponse>(`${apiUrl}/api/videos`)
      .then((response: AxiosResponse<VideoResponse>) => {
        setVideos(response.data.videos);
      })
      .catch((error: unknown) => {
        console.error('There was an error fetching the videos!', error);
      });
  }, []);

  return (
    <div className="home">
      <h1>🎬 Watch UA</h1>
      <p>Discover amazing video content</p>
      <div className="video-list">
        {videos.map(video => (
          <div key={video.id} className="video-card" onClick={() => window.location.href = `/video/${video.id}`}>
            <img src={video.previewImage} alt={video.title} />
            <div className="video-content">
              <h2>{video.title}</h2>
              <p>{video.description}</p>
              <div className="video-stats">
                <span>👍 {video.likes}</span>
                <span>👎 {video.dislikes}</span>
                <span>💬 {video.comments}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Home;
