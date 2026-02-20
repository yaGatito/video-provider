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
    axios.get<VideoResponse>('/api/videos')
      .then((response: AxiosResponse<VideoResponse>) => {
        setVideos(response.data.videos);
      })
      .catch((error: unknown) => {
        console.error('There was an error fetching the videos!', error);
      });
  }, []);

  return (
    <div className="home">
      <h1>Watch UA</h1>
      <div className="video-list">
        {videos.map(video => (
          <div key={video.id} className="video-card" onClick={() => window.location.href = `/video/${video.id}`}>
            <img src={video.previewImage} alt={video.title} />
            <h2>{video.title}</h2>
            <p>{video.description}</p>
            <div className="video-stats">
              <span>{video.likes} Likes</span>
              <span>{video.dislikes} Dislikes</span>
              <span>{video.comments} Comments</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Home;
