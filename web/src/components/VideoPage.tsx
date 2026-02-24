import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import { useParams } from 'react-router-dom';
import './VideoPage.css';

interface Video {
    id: number;
    title: string;
    description: string;
    likes: number;
    dislikes: number;
    comments: number;
    previewImage: string;
}

const VideoPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [video, setVideo] = useState<Video | null>(null);

    useEffect(() => {
        const apiUrl = process.env.REACT_APP_API_URL;
        axios.get<Video>(`${apiUrl}/v1/videos/id/${id}`)
            .then((response: AxiosResponse<Video>) => {
                setVideo(response.data);
            })
            .catch((error: unknown) => {
                console.error('There was an error fetching the video!', error);
            });
    }, [id]);

    if (!video) return <div className="loading">Loading video...</div>;

    return (
        <div className="video-page">
            <div className="video-page-container">
                <div className="video-page-header">
                    <a href="/" className="back-button">← Back to Videos</a>
                    <h1>{video.title}</h1>
                </div>
                <div className="video-page-content">
                    <img src={video.previewImage} alt={video.title} />
                    <p>{video.description}</p>
                    <div className="video-stats">
                        <span>👍 {video.likes} Likes</span>
                        <span>👎 {video.dislikes} Dislikes</span>
                        <span>💬 {video.comments} Comments</span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default VideoPage;