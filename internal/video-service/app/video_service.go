package app

import (
	"context"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/ports"
)

type VideoService interface {

	// Create creates a new video entity in the storage.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - video: The video data to be created.
	//
	// Returns:
	// - domain.Video: The created video.
	// - error: An error if the video could not be created.
	Create(ctx context.Context,
		video domain.Video,
	) (domain.Video, error)

	// GetByID retrieves a video by its unique identifier.
	//
	// Parameters:
	// - ctx: The context for the operation, which may include cancellation signals or deadlines.
	// - videoID: The unique identifier of the video to retrieve.
	//
	// Returns:
	// - domain.Video: The retrieved video if found.
	// - error: An error if the video is not found or if an unexpected issue occurs.
	GetByID(ctx context.Context,
		videoID domain.UUID,
	) (domain.Video, error)

	// GetByPublisher retrieves a list of videos by a specific publisher.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - publisherID: The unique identifier of the publisher.
	// - params: Pagination and sorting parameters.
	//
	// Returns:
	// - []domain.Video: A list of videos.
	// - error: An error if the operation fails.
	GetByPublisher(ctx context.Context,
		publisherID domain.UUID,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	// SearchPublisher searches for videos by a specific publisher.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - publisherID: The unique identifier of the publisher.
	// - query: A search query.
	// - params: Pagination and sorting parameters.
	//
	// Returns:
	// - []domain.Video: A list of videos matching the query.
	// - error: An error if the operation fails.
	SearchPublisher(ctx context.Context,
		publisherID domain.UUID,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)

	// SearchGlobal searches for videos based on a global query.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - query: The search term.
	// - params: Pagination and sorting parameters.
	//
	// Returns:
	// - []domain.Video: A list of videos matching the query.
	// - error: An error if the operation fails.
	SearchGlobal(ctx context.Context,
		query string,
		params domain.VideoPageParams,
	) ([]domain.Video, error)
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

func (vs *VideoInteractor) Create(ctx context.Context,
	video domain.Video,
) (domain.Video, error) {
	return vs.repo.CreateVideo(ctx, video)
}

func (vs *VideoInteractor) GetByID(ctx context.Context,
	videoID domain.UUID,
) (domain.Video, error) {
	return vs.repo.GetVideoByID(ctx, videoID)
}

func (vs *VideoInteractor) GetByPublisher(ctx context.Context,
	publisherID domain.UUID,
	params domain.VideoPageParams,
) ([]domain.Video, error) {
	return vs.repo.GetPublisherVideos(ctx, publisherID, params)
}

func (s *VideoInteractor) SearchPublisher(ctx context.Context,
	publisherID domain.UUID,
	query string,
	params domain.VideoPageParams,
) ([]domain.Video, error) {
	return s.repo.SearchPublisher(ctx, publisherID, query, params)
}

func (s *VideoInteractor) SearchGlobal(ctx context.Context,
	query string, params domain.VideoPageParams,
) ([]domain.Video, error) {
	return s.repo.SearchGlobal(ctx, query, params)
}
