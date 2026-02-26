package app

import (
	"context"
	"fmt"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/ports"
)

type VideoService interface {
	// Create creates a new video in the storage.
	//
	// Parameters:
	// - ctx: The context for the operation, which may include cancellation signals or deadlines.
	// - video: The video data to be created, which must include valid metadata such as title, description.
	//
	// Returns:
	// - domain.Video: The created video, including any system-generated fields like ID or timestamps.
	// - error: An error if the video could not be created, such as due to invalid input or duplicate data.
	//
	// Validation Rules:
	// - The video's title must not be empty.
	// - The video's description must not be empty.
	// - The video's URL must be a valid, non-empty string.
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
	//
	// Validation Rules:
	// - The videoID must be a valid, non-empty UUID.
	GetByID(ctx context.Context,
		videoID domain.UUID,
	) (domain.Video, error)

	// GetByPublisher retrieves a list of videos associated with a specific publisher.
	// It allows sorting the results by a specified field, applying pagination via offset and limit,
	// and defining the sort order (ascending or descending).
	//
	// Parameters:
	//   - ctx: The context for the operation, which may carry deadlines, cancellation signals, and other request-scoped values.
	//   - publisherID: The unique identifier of the publisher whose videos are to be retrieved.
	//   - orderBy: The field by which the results should be sorted. Valid values are determined by the implementation.
	//   - offset: The number of records to skip. Used for pagination.
	//   - limit: The maximum number of records to return. Used for pagination.
	//   - asc: The sort order. If "true", the results are sorted in ascending order; otherwise, in descending order.
	//
	// Returns:
	//   - []domain.Video: A slice of video objects retrieved based on the given criteria.
	//   - error: An error if the operation fails, such as due to invalid input or database issues.
	//
	// Validation Rules:
	//   - publisherID must be a valid UUID.
	//   - orderBy must be a valid field name supported by the sorting mechanism.
	//   - offset and limit must be non-negative integers.
	//   - asc must be either "true" or "false".
	GetByPublisher(ctx context.Context,
		publisherID domain.UUID,
		orderBy string,
		offset int32,
		limit int32,
		asc string,
	) ([]domain.Video, error)

	// SearchPublisher searches for videos by a specific publisher.
	// It allows filtering by query, sorting by a specific field, and pagination.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - publisherID: The unique identifier of the publisher to search videos for.
	// - query: A search query string to filter videos (e.g., title, description).
	// - orderBy: The field to sort the results by (e.g., "title", "date").
	// - offset: The number of records to skip (used for pagination).
	// - limit: The maximum number of records to return (used for pagination).
	// - asc: The sorting order, either "asc" for ascending or "desc" for descending.
	//
	// Returns:
	// - A slice of Video objects that match the search criteria.
	// - An error if the operation fails.
	//
	// Validation Rules:
	// - publisherID must be a valid UUID.
	// - orderBy must be a valid field name supported by the video model.
	// - offset and limit must be non-negative integers.
	// - asc must be either "asc" or "desc".
	SearchPublisher(ctx context.Context,
		publisherID domain.UUID,
		query string,
		orderBy string,
		offset int32,
		limit int32,
		asc string,
	) ([]domain.Video, error)

	// SearchGlobal searches for videos based on a global query.
	// It allows sorting the results by a specified field, with options for ascending or descending order,
	// and supports pagination using offset and limit parameters.
	//
	// Parameters:
	//   - ctx: The context for the function call, which may carry deadlines, cancellation signals, and other request-scoped values.
	//   - query: The search term used to filter videos. This can be any string that matches video metadata.
	//   - orderBy: The field by which the results should be sorted. Valid fields include "title", "date", "views", etc.
	//   - offset: The number of records to skip before starting to return results. Used for pagination.
	//   - limit: The maximum number of results to return. Used for pagination.
	//   - asc: A string indicating whether the sorting should be in ascending or descending order. Valid values are "asc" or "desc".
	//
	// Returns:
	//   - []domain.Video: A slice of video objects that match the search query.
	//   - error: An error if the search fails for any reason, such as invalid parameters or database errors.
	//
	// Validation:
	//   - The `query` parameter must be a non-empty string.
	//   - The `orderBy` parameter must be a valid field name supported by the video database.
	//   - The `offset` and `limit` parameters must be non-negative integers.
	//   - The `asc` parameter must be either "t" or "f".
	//
	// Example:
	//   videos, err := videoService.SearchGlobal(ctx, "golang", "createdAt", 0, 10, "t")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   fmt.Println(videos)
	SearchGlobal(ctx context.Context,
		query string,
		orderBy string,
		offset int32,
		limit int32,
		asc string,
	) ([]domain.Video, error)
}

type VideoInteractor struct {
	repo ports.VideoRepository
}

var nilUUID domain.UUID
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

	if err := video.Validate(); err != nil {
		return domain.Video{}, err
	}
	return vs.repo.CreateVideo(ctx, video)
}

func (vs *VideoInteractor) GetByID(ctx context.Context,
	videoID domain.UUID,
) (domain.Video, error) {

	if videoID == nilUUID {
		return domain.Video{}, fmt.Errorf("empty video ID")
	}

	return vs.repo.GetVideoByID(ctx, videoID)
}

func (vs *VideoInteractor) GetByPublisher(ctx context.Context,
	publisherID domain.UUID,
	orderBy string,
	offset int32,
	limit int32,
	asc string,
) ([]domain.Video, error) {

	if publisherID == nilUUID {
		return nil, fmt.Errorf("empty publisher ID")
	}

	pageParams, err := ValidVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		return nil, err
	}

	return vs.repo.GetPublisherVideos(ctx, publisherID, pageParams)
}

func (s *VideoInteractor) SearchPublisher(ctx context.Context,
	publisherID domain.UUID,
	query string,
	orderBy string,
	offset int32,
	limit int32,
	asc string,
) ([]domain.Video, error) {

	if publisherID == nilUUID {
		return nil, fmt.Errorf("empty publisher ID")
	}
	query, err := ValidateSearchQuery(query)
	if err != nil {
		return nil, err
	}

	pageParams, err := ValidVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		return nil, err
	}

	return s.repo.SearchPublisher(ctx,
		publisherID,
		ports.VideoSearchParams{
			Query:           query,
			VideoPageParams: pageParams,
		})
}

func (s *VideoInteractor) SearchGlobal(ctx context.Context,
	query string,
	orderBy string,
	offset int32,
	limit int32,
	asc string,
) ([]domain.Video, error) {

	query, err := ValidateSearchQuery(query)
	if err != nil {
		return nil, err
	}

	pageParams, err := ValidVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		return nil, err
	}

	return s.repo.SearchGlobal(ctx,
		ports.VideoSearchParams{
			Query:           query,
			VideoPageParams: pageParams,
		})
}

func ValidVideoPageParams(
	orderBy string,
	offset int32,
	limit int32,
	asc string,
) (ports.VideoPageParams, error) {

	offset, err := ValidateOffset(offset)
	if err != nil {
		return ports.VideoPageParams{}, err
	}
	limit, err = ValidateLimit(limit)
	if err != nil {
		return ports.VideoPageParams{}, err
	}
	orderBy, err = ValidateOrderBy(orderBy)
	if err != nil {
		return ports.VideoPageParams{}, err
	}
	isAsc, err := ValidateIsAsc(asc)
	if err != nil {
		return ports.VideoPageParams{}, err
	}

	return ports.VideoPageParams{
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
		IsAsc:   isAsc,
	}, nil
}
