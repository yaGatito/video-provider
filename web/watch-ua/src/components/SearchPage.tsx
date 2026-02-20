import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './SearchPage.css';

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

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    if (query) {
      const apiUrl = process.env.REACT_APP_API_URL;
      axios.get<VideoResponse>(`${apiUrl}/api/videos?search=${query}`)
        .then((response: AxiosResponse<VideoResponse>) => {
          setVideos(response.data.videos);
        })
        .catch((error: unknown) => {
          console.error('There was an error fetching the search results!', error);
        });
    }
  }, [query]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
  };

  return (
    <div className="search-page">
      <h1>🔍 Search</h1>
      <input
        type="text"
        placeholder="Search for videos..."
        value={query}
        onChange={handleSearch}
      />
      {videos.length > 0 ? (
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
      ) : (
        <p>No videos found.</p>
      )}
    </div>
  );
};

export default SearchPage;