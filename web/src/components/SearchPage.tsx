import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import './SearchPage.css';

interface Video {
  id: string;
  publisherID: string;
  topic: string;
  description: string;
  createdAt: Date;
}

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [limit, setLimit] = useState(10);
  const [offset, setOffset] = useState(0);
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    if (query.length > 2) {
      const apiUrl = process.env.REACT_APP_API_URL
      axios.get<Video[]>(`${apiUrl}/v1/videos/search?query=${query}&limit=${limit}&offset=${offset}`)
        .then((response: AxiosResponse<Video[]>) => {
          setVideos(response.data);
        })
        .catch((error: unknown) => {
          console.error('There was an error fetching the search results!', error);
        });
    }
  }, [query, limit, offset]);

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
    setOffset(0); // Reset offset when query changes
  };

  const loadMoreVideos = () => {
    setOffset(offset + limit);
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
        <>
          <div className="video-list">
            {videos.map(video => (
              <div key={video.id} className="video-card" onClick={() => window.location.href = `/v1/videos/${video.id}`}>
                <div className="video-content">
                  <h2>{video.topic}</h2>
                  <p>{video.description}</p>
                  <div className="video-stats">
                    <span>Published on {new Date(video.createdAt).toLocaleDateString()}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <button onClick={loadMoreVideos}>Load More</button>
        </>
      ) : (
        <p>No videos found.</p>
      )}
    </div>
  );
};

export default SearchPage;