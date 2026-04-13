import React from 'react';
import styled from 'styled-components';
import Button from './Button';

export interface VideoPreview {
  id: string;
  topic: string;
  description: string;
  createdAt?: string;
  thumbnailUrl?: string;
  comments?: number;
}

interface VideoCardProps {
  video: VideoPreview;
  onClick?: () => void;
  onComment?: () => void;
  onSubscribe?: () => void;
}

const Card = styled.article<{ $clickable?: boolean }>`
  display: grid;
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  overflow: hidden;
  box-shadow: ${({ theme }) => theme.shadows.sm};
  transition: transform ${({ theme }) => theme.transitions.base},
    box-shadow ${({ theme }) => theme.transitions.base};
  cursor: ${({ $clickable }) => ($clickable ? 'pointer' : 'default')};

  &:hover {
    transform: ${({ $clickable }) => ($clickable ? 'translateY(-4px)' : 'none')};
    box-shadow: ${({ theme, $clickable }) => ($clickable ? theme.shadows.md : theme.shadows.sm)};
  }
`;

const Preview = styled.img`
  width: 100%;
  height: 200px;
  object-fit: cover;
  display: block;
  background: ${({ theme }) => theme.colors.surfaceAlt};
`;

const CardBody = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
  padding: ${({ theme }) => theme.spacing.lg};
`;

const Title = styled.h3`
  font-family: ${({ theme }) => theme.fonts.heading};
  margin: 0;
`;

const Description = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin: 0;
  min-height: 3.6rem;
`;

const Meta = styled.div`
  font-size: 0.9rem;
  color: ${({ theme }) => theme.colors.muted};
`;

const ActionRow = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
  margin-top: ${({ theme }) => theme.spacing.md};
  justify-items: start;
`;

const VideoCard: React.FC<VideoCardProps> = ({ video, onClick, onComment, onSubscribe }) => {
  const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, callback?: () => void) => {
    event.stopPropagation();
    callback?.();
  };

  return (
    <Card $clickable={Boolean(onClick)} onClick={onClick}>
      <Preview
        src={video.thumbnailUrl || 'https://images.unsplash.com/photo-1618345522246-81e1f1e041b7?auto=format&fit=crop&w=1920&q=80'}
        alt={video.topic}
      />
      <CardBody>
        <Title>{video.topic}</Title>
        <Description>{video.description}</Description>
        <Meta>{video.createdAt ? `Published ${new Date(video.createdAt).toLocaleDateString()}` : 'Published recently'}</Meta>

        {(onComment || onSubscribe) && (
          <ActionRow>
            {onSubscribe && (
              <Button
                variant="secondary"
                onClick={(event) => handleButtonClick(event, onSubscribe)}
              >
                Subscribe
              </Button>
            )}
            {onComment && (
              <Button
                variant="ghost"
                onClick={(event) => handleButtonClick(event, onComment)}
              >
                Comment
              </Button>
            )}
          </ActionRow>
        )}
      </CardBody>
    </Card>
  );
};

export default VideoCard;
