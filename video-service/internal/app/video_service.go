package app

import (
	"context"
	"fmt"
	"video-service/internal/domain"
	"video-service/internal/ports"
)

// mockgen -source="./internal/app/service.go" -destination="./internal/app/mock/service.go" -mock_names=("VideoInteractor")

type VideoService interface {
	Create(ctx context.Context, video domain.Video) error
	GetByID(ctx context.Context, videoID domain.UUID) (domain.Video, error)
	GetByPublisher(ctx context.Context, publisherID domain.UUID, offset, limit int32) ([]domain.Video, error)
	SearchPublisher(ctx context.Context, publisherID domain.UUID, query string, offset, limit int32) ([]domain.Video, error)
	SearchGlobal(ctx context.Context, query string, offset, limit int32) ([]domain.Video, error)
}

type VideoInteractor struct {
	repo ports.VideoRepository
}

var _ VideoService = (*VideoInteractor)(nil)

type VideosResultList struct {
	Videos []domain.Video
	Size   int
}

func NewVideoInteractor(repo ports.VideoRepository) VideoService {
	return &VideoInteractor{repo: repo}
}

func (vs *VideoInteractor) Create(ctx context.Context, video domain.Video) error {
	if err := video.Validate(); err != nil {
		return err
	}
	return vs.repo.CreateVideo(ctx, video)
}

var nilUUID = domain.UUID{}

func (vs *VideoInteractor) GetByID(ctx context.Context, videoID domain.UUID) (domain.Video, error) {
	if videoID == nilUUID {
		return domain.Video{}, fmt.Errorf("empty video ID")
	}

	return vs.repo.GetVideoByID(ctx, videoID)
}

func (vs *VideoInteractor) GetByPublisher(ctx context.Context, publisherID domain.UUID, offset, limit int32) ([]domain.Video, error) {
	if publisherID == nilUUID {
		return nil, fmt.Errorf("empty publisher ID")
	}
	offset, limit = ValidatePagination(offset, limit)

	return vs.repo.GetPublisherVideos(ctx, publisherID, ports.PageRequest{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
}

func (s *VideoInteractor) SearchPublisher(ctx context.Context, publisherID domain.UUID, query string, offset, limit int32) ([]domain.Video, error) {
	if publisherID == nilUUID {
		return nil, fmt.Errorf("empty publisher ID")
	}
	query, err := ValidateSearchQuery(query)
	if err != nil {
		return nil, err
	}
	offset, limit = ValidatePagination(offset, limit)

	return s.repo.SearchPublisher(ctx, publisherID, ports.VideoSearch{
		Query: query,
		PageRequest: ports.PageRequest{
			Offset: int32(offset),
			Limit:  int32(limit),
		}})
}

func (s *VideoInteractor) SearchGlobal(ctx context.Context, query string, offset, limit int32) ([]domain.Video, error) {
	query, err := ValidateSearchQuery(query)
	if err != nil {
		return nil, err
	}
	offset, limit = ValidatePagination(offset, limit)

	return s.repo.SearchGlobal(ctx, ports.VideoSearch{
		Query: query,
		PageRequest: ports.PageRequest{
			Offset: int32(offset),
			Limit:  int32(limit),
		}})
}
