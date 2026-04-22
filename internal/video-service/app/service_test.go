package app_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"video-provider/video-service/app"
	"video-provider/video-service/domain"
	mock_ports "video-provider/video-service/ports/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoService_Create(t *testing.T) {
	testCases := []struct {
		testName string
		video    domain.Video
		expected domain.Video
		err      error
	}{
		{
			testName: "Valid video creation",
			video: domain.Video{
				ID:          uuid.New(),
				PublisherID: uuid.New(),
				Topic:       "Amazing Travel Vlog",
				Description: "This is a travel vlog about exploring new countries",
				CreatedAt:   time.Now(),
				Status:      domain.StatusDraft,
			},
			expected: domain.Video{
				ID:          uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
				PublisherID: uuid.MustParse("a9fa522f-0006-464f-8d68-356ba1d6ad7d"),
				Topic:       "Amazing Travel Vlog",
				Description: "This is a travel vlog about exploring new countries",
				CreatedAt:   time.Now(),
				Status:      domain.StatusDraft,
			},
			err: nil,
		},
		{
			testName: "Database error during creation",
			video: domain.Video{
				PublisherID: uuid.New(),
				Topic:       "Test Video",
				Description: "Test Description",
				Status:      domain.StatusDraft,
			},
			expected: domain.Video{},
			err:      errors.New("database error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockVideoRepository(ctrl)
			videoService := app.NewVideoInteractor(mockRepo)

			mockRepo.EXPECT().
				CreateVideo(gomock.Any(), gomock.Any()).
				Return(tc.expected, tc.err).
				Times(1)

			result, err := videoService.Create(context.Background(), tc.video)

			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expected, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Topic, result.Topic)
				assert.Equal(t, tc.expected.PublisherID, result.PublisherID)
			}
		})
	}
}

func TestVideoService_GetByID(t *testing.T) {
	expectedVideoID := uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	expectedPublisherID := uuid.MustParse("a9fa522f-0006-464f-8d68-356ba1d6ad7d")

	testCases := []struct {
		testName string
		videoID  uuid.UUID
		expected domain.Video
		err      error
	}{
		{
			testName: "Valid video retrieval",
			videoID:  expectedVideoID,
			expected: domain.Video{
				ID:          expectedVideoID,
				PublisherID: expectedPublisherID,
				Topic:       "Great Video",
				Description: "This is a great video",
				CreatedAt:   time.Now(),
				Status:      domain.StatusPublished,
			},
			err: nil,
		},
		{
			testName: "Video not found",
			videoID:  uuid.New(),
			expected: domain.Video{},
			err:      errors.New("video not found"),
		},
		{
			testName: "Database error",
			videoID:  uuid.New(),
			expected: domain.Video{},
			err:      errors.New("database error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockVideoRepository(ctrl)
			videoService := app.NewVideoInteractor(mockRepo)

			mockRepo.EXPECT().
				GetVideoByID(gomock.Any(), tc.videoID).
				Return(tc.expected, tc.err).
				Times(1)

			result, err := videoService.GetByID(context.Background(), tc.videoID)

			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expected, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Topic, result.Topic)
			}
		})
	}
}

func TestVideoService_GetByPublisher(t *testing.T) {
	publisherID := uuid.MustParse("a9fa522f-0006-464f-8d68-356ba1d6ad7d")
	videoID1 := uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	videoID2 := uuid.MustParse("e9fa522f-0006-464f-8d68-356ba1d6ad7d")

	testCases := []struct {
		testName    string
		publisherID uuid.UUID
		params      domain.VideoPageParams
		expected    []domain.Video
		err         error
	}{
		{
			testName:    "Get publisher videos with valid params",
			publisherID: publisherID,
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{
				{
					ID:          videoID1,
					PublisherID: publisherID,
					Topic:       "Video 1",
					Description: "First video",
					CreatedAt:   time.Now(),
					Status:      domain.StatusPublished,
				},
				{
					ID:          videoID2,
					PublisherID: publisherID,
					Topic:       "Video 2",
					Description: "Second video",
					CreatedAt:   time.Now(),
					Status:      domain.StatusPublished,
				},
			},
			err: nil,
		},
		{
			testName:    "Get publisher videos - no videos",
			publisherID: uuid.New(),
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{},
			err:      nil,
		},
		{
			testName:    "Database error",
			publisherID: publisherID,
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: nil,
			err:      errors.New("database error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockVideoRepository(ctrl)
			videoService := app.NewVideoInteractor(mockRepo)

			mockRepo.EXPECT().
				GetPublisherVideos(gomock.Any(), tc.publisherID, tc.params).
				Return(tc.expected, tc.err).
				Times(1)

			result, err := videoService.GetByPublisher(context.Background(), tc.publisherID, tc.params)

			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tc.expected), len(result))
				if len(tc.expected) > 0 {
					assert.Equal(t, tc.expected[0].ID, result[0].ID)
				}
			}
		})
	}
}

func TestVideoService_SearchPublisher(t *testing.T) {
	publisherID := uuid.MustParse("a9fa522f-0006-464f-8d68-356ba1d6ad7d")
	videoID := uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")

	testCases := []struct {
		testName    string
		publisherID uuid.UUID
		query       string
		params      domain.VideoPageParams
		expected    []domain.Video
		err         error
	}{
		{
			testName:    "Valid publisher search",
			publisherID: publisherID,
			query:       "travel",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{
				{
					ID:          videoID,
					PublisherID: publisherID,
					Topic:       "Travel Vlog",
					Description: "Exploring the world",
					CreatedAt:   time.Now(),
					Status:      domain.StatusPublished,
				},
			},
			err: nil,
		},
		{
			testName:    "No results for search",
			publisherID: publisherID,
			query:       "nonexistent",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{},
			err:      nil,
		},
		{
			testName:    "Database error on search",
			publisherID: publisherID,
			query:       "test",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: nil,
			err:      errors.New("database error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockVideoRepository(ctrl)
			videoService := app.NewVideoInteractor(mockRepo)

			mockRepo.EXPECT().
				SearchPublisher(gomock.Any(), tc.publisherID, tc.query, tc.params).
				Return(tc.expected, tc.err).
				Times(1)

			result, err := videoService.SearchPublisher(context.Background(), tc.publisherID, tc.query, tc.params)

			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tc.expected), len(result))
			}
		})
	}
}

func TestVideoService_SearchGlobal(t *testing.T) {
	videoID1 := uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	videoID2 := uuid.MustParse("e9fa522f-0006-464f-8d68-356ba1d6ad7d")
	publisherID1 := uuid.MustParse("a9fa522f-0006-464f-8d68-356ba1d6ad7d")
	publisherID2 := uuid.MustParse("b9fa522f-0006-464f-8d68-356ba1d6ad7d")

	testCases := []struct {
		testName string
		query    string
		params   domain.VideoPageParams
		expected []domain.Video
		err      error
	}{
		{
			testName: "Valid global search",
			query:    "tutorial",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{
				{
					ID:          videoID1,
					PublisherID: publisherID1,
					Topic:       "Python Tutorial",
					Description: "Learn Python basics",
					CreatedAt:   time.Now(),
					Status:      domain.StatusPublished,
				},
				{
					ID:          videoID2,
					PublisherID: publisherID2,
					Topic:       "Go Tutorial",
					Description: "Learn Go programming",
					CreatedAt:   time.Now(),
					Status:      domain.StatusPublished,
				},
			},
			err: nil,
		},
		{
			testName: "No results from global search",
			query:    "xyz",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: []domain.Video{},
			err:      nil,
		},
		{
			testName: "Database error on global search",
			query:    "test",
			params: domain.VideoPageParams{
				OrderBy: domain.OrderByDate,
				Offset:  0,
				Limit:   10,
				Asc:     domain.DescOrder,
			},
			expected: nil,
			err:      errors.New("database error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockVideoRepository(ctrl)
			videoService := app.NewVideoInteractor(mockRepo)

			mockRepo.EXPECT().
				SearchGlobal(gomock.Any(), tc.query, tc.params).
				Return(tc.expected, tc.err).
				Times(1)

			result, err := videoService.SearchGlobal(context.Background(), tc.query, tc.params)

			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tc.expected), len(result))
			}
		})
	}
}
