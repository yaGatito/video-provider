package ports

import (
	"context"
	"video-provider/video-service/domain"

	"github.com/google/uuid"
)

type VideoRepository interface {
	CreateVideo(
		ctx context.Context,
		video domain.Video,
	) (domain.Video, error)

	GetVideoByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Video, error)

	GetPublisherVideos(ctx context.Context,
		publisherID uuid.UUID,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	SearchPublisher(
		ctx context.Context,
		publisherID uuid.UUID,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	SearchGlobal(
		ctx context.Context,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)
}
