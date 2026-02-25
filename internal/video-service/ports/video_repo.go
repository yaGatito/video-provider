package ports

import (
	"context"
	"video-provider/internal/video-service/domain"
)

type VideoPageParams struct {
	Offset int32
	Limit  int32
	SortBy string
	IsAsc  bool
}

type VideoSearchParams struct {
	VideoPageParams
	Query string
}

const (
	CreatedAtSort string = "createdAt"
	ascOrder      string = "ASC"
	descOrder     string = "DESC"
)

type VideoRepository interface {
	CreateVideo(ctx context.Context, video domain.Video) (domain.Video, error)
	GetVideoByID(ctx context.Context, id domain.UUID) (domain.Video, error)
	GetPublisherVideos(
		ctx context.Context,
		publisherID domain.UUID,
		args VideoPageParams,
	) ([]domain.Video, error)
	SearchPublisher(
		ctx context.Context,
		publisherID domain.UUID,
		args VideoSearchParams,
	) ([]domain.Video, error)
	SearchGlobal(ctx context.Context, args VideoSearchParams) ([]domain.Video, error)
}
