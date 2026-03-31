package ports

import (
	"context"
	"video-provider/internal/video-service/domain"
)

type VideoRepository interface {
	CreateVideo(
		ctx context.Context,
		video domain.Video,
	) (domain.Video, error)

	GetVideoByID(
		ctx context.Context,
		id domain.UUID,
	) (domain.Video, error)

	GetPublisherVideos(ctx context.Context,
		publisherID domain.UUID,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	SearchPublisher(
		ctx context.Context,
		publisherID domain.UUID,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	SearchGlobal(
		ctx context.Context,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)
}
