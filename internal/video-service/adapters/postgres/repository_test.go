package postgres_test

import (
	"context"
	"testing"
	"time"

	"video-provider/video-service/adapters/postgres"
	"video-provider/video-service/adapters/postgres/sqlcgen"
	"video-provider/video-service/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateVideo(
	ctx context.Context,
	params sqlcgen.CreateVideoParams,
) (sqlcgen.Video, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sqlcgen.Video), args.Error(1)
}

func (m *MockQuerier) GetVideoByID(ctx context.Context, id uuid.UUID) (sqlcgen.Video, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlcgen.Video), args.Error(1)
}

func (m *MockQuerier) GetVideosByPublisher(
	ctx context.Context,
	params sqlcgen.GetVideosByPublisherParams,
) ([]sqlcgen.Video, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sqlcgen.Video), args.Error(1)
}

func (m *MockQuerier) SearchPublisher(
	ctx context.Context,
	params sqlcgen.SearchPublisherParams,
) ([]sqlcgen.Video, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sqlcgen.Video), args.Error(1)
}

func (m *MockQuerier) SearchGlobal(
	ctx context.Context,
	params sqlcgen.SearchGlobalParams,
) ([]sqlcgen.Video, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sqlcgen.Video), args.Error(1)
}

func TestVideoRepoPostgreSQL_CreateVideo(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewVideoRepoPostgreSQL(mockQuerier)

	video := domain.Video{
		PublisherID: uuid.New(),
		Topic:       "Test Topic",
		Description: "Test Description",
	}

	mockVideo := sqlcgen.Video{
		ID:          uuid.New(),
		Publisherid: video.PublisherID,
		Topic:       video.Topic,
		Description: pgtype.Text{String: video.Description, Valid: true},
		Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
	}
	mockQuerier.On("CreateVideo", mock.Anything, sqlcgen.CreateVideoParams{
		Publisherid: video.PublisherID,
		Topic:       video.Topic,
		Description: pgtype.Text{String: video.Description, Valid: true},
	}).Return(mockVideo, nil)

	result, err := repo.CreateVideo(context.Background(), video)
	assert.NoError(t, err)
	assert.Equal(t, mockVideo.ID, result.ID)
}

func TestVideoRepoPostgreSQL_GetVideoByID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewVideoRepoPostgreSQL(mockQuerier)

	videoID := uuid.New()

	mockVideo := sqlcgen.Video{
		ID:          videoID,
		Publisherid: uuid.New(),
		Topic:       "Test Topic",
		Description: pgtype.Text{String: "Test Description", Valid: true},
		Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
	}
	mockQuerier.On("GetVideoByID", mock.Anything, videoID).Return(mockVideo, nil)

	video, err := repo.GetVideoByID(context.Background(), videoID)
	assert.NoError(t, err)
	assert.Equal(t, videoID, video.ID)
}

func TestVideoRepoPostgreSQL_GetPublisherVideos(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewVideoRepoPostgreSQL(mockQuerier)

	publisherID := uuid.New()
	params := domain.VideoPageParams{
		Offset: 0,
		Limit:  10,
	}

	mockVideos := []sqlcgen.Video{
		{
			ID:          uuid.New(),
			Publisherid: publisherID,
			Topic:       "Test Topic 1",
			Description: pgtype.Text{String: "Test Description 1", Valid: true},
			Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
		},
		{
			ID:          uuid.New(),
			Publisherid: publisherID,
			Topic:       "Test Topic 2",
			Description: pgtype.Text{String: "Test Description 2", Valid: true},
			Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
		},
	}
	mockQuerier.On("GetVideosByPublisher", mock.Anything, sqlcgen.GetVideosByPublisherParams{
		Publisherid: publisherID,
		Offset:      params.Offset,
		Limit:       params.Limit,
	}).Return(mockVideos, nil)

	videos, err := repo.GetPublisherVideos(context.Background(), publisherID, params)
	assert.NoError(t, err)
	assert.Len(t, videos, 2)
}

func TestVideoRepoPostgreSQL_SearchPublisher(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewVideoRepoPostgreSQL(mockQuerier)

	publisherID := uuid.New()
	query := "Test Query"
	params := domain.VideoPageParams{
		Offset: 0,
		Limit:  10,
	}

	mockVideos := []sqlcgen.Video{
		{
			ID:          uuid.New(),
			Publisherid: publisherID,
			Topic:       "Test Topic 1",
			Description: pgtype.Text{String: "Test Description 1", Valid: true},
			Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
		},
	}
	mockQuerier.On("SearchPublisher", mock.Anything, sqlcgen.SearchPublisherParams{
		Publisherid: publisherID,
		Column2:     query,
		Column3:     postgres.GetOrderBy(params.OrderBy, params.Asc),
		Offset:      params.Offset,
		Limit:       params.Limit,
	}).Return(mockVideos, nil)

	videos, err := repo.SearchPublisher(context.Background(), publisherID, query, params)
	assert.NoError(t, err)
	assert.Len(t, videos, 1)
}

func TestVideoRepoPostgreSQL_SearchGlobal(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewVideoRepoPostgreSQL(mockQuerier)

	query := "Test Query"
	params := domain.VideoPageParams{
		Offset: 0,
		Limit:  10,
	}

	mockVideos := []sqlcgen.Video{
		{
			ID:          uuid.New(),
			Publisherid: uuid.New(),
			Topic:       "Test Topic 1",
			Description: pgtype.Text{String: "Test Description 1", Valid: true},
			Createdat:   pgtype.Time{Microseconds: time.Now().UnixMicro(), Valid: true},
		},
	}
	mockQuerier.On("SearchGlobal", mock.Anything, sqlcgen.SearchGlobalParams{
		Column1: query,
		Column2: postgres.GetOrderBy(params.OrderBy, params.Asc),
		Offset:  params.Offset,
		Limit:   params.Limit,
	}).Return(mockVideos, nil)

	videos, err := repo.SearchGlobal(context.Background(), query, params)
	assert.NoError(t, err)
	assert.Len(t, videos, 1)
}
