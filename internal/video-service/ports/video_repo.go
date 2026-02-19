package ports

import (
	"context"
	"video-provider/internal/video-service/domain"
)

type VideoRepository interface {
	CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	GetVideoByID(ctx context.Context, id domain.UUID) (domain.Video, error)
	GetPublisherVideos(
		ctx context.Context,
		publisherID domain.UUID,
		args PageRequest,
	) ([]domain.Video, error)
	SearchPublisher(
		ctx context.Context,
		publisherID domain.UUID,
		args VideoSearch,
	) ([]domain.Video, error)
	SearchGlobal(ctx context.Context, args VideoSearch) ([]domain.Video, error)
}

type PageRequest struct {
	Offset int32
	Limit  int32
}

type VideoSearch struct {
	PageRequest
	Query string
}
