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
        axios.get<Video>(`/api/videos/${id}`)
            .then((response: AxiosResponse<Video>) => {
                setVideo(response.data);
            })
            .catch((error: unknown) => {
                console.error('There was an error fetching the video!', error);
            });
    }, [id]);

    if (!video) return <div>Loading...</div>;

    return (
        <div className="video-page">
            <h1>{video.title}</h1>
            <img src={video.previewImage} alt={video.title} />
            <p>{video.description}</p>
            <div className="video-stats">
                <span>{video.likes} Likes</span>
                <span>{video.dislikes} Dislikes</span>
                <span>{video.comments} Comments</span>
            </div>
        </div>
    );
};

export default VideoPage;