package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"video-provider/internal/video-service/app"
	"video-provider/internal/video-service/domain"
	mock_ports "video-provider/internal/video-service/ports/mock"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

var testTopic = "topic"
var testDesc = "desc"
var testPublisherID, _ = uuid.Parse("d9fa522f-0006-464f-8d68-326ba1d6ad7d")
var testVideo = domain.Video{
	PublisherID: testPublisherID,
	Topic:       testTopic,
	Description: testDesc,
}
var testVideoID, _ = uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")

func TestCreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	repo.
		EXPECT().
		CreateVideo(gomock.Any(), gomock.Eq(testVideo)).
		MaxTimes(1)

	_, err := videoService.Create(context.Background(), testVideo)
	require.NoError(t, err)
}

func TestCreateInvalidVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)
	invalidVideo := domain.Video{}

	repo.
		EXPECT().
		CreateVideo(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	res, err := videoService.Create(context.Background(), invalidVideo)
	require.Error(t, err, "expected error for invalid video"+err.Error())
	require.NotNil(t, res, fmt.Sprintf("expected empty video for invalid video: %v", res))
}

func TestGetdVideoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	expectedVideo := domain.Video{
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: testDesc,
		CreatedAt:   time.Now(),
	}

	cases := []struct {
		name    string
		videoID uuid.UUID
		wantErr bool
	}{
		{"ok", testVideoID, false},
		{"no video id", uuid.Nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				// Error scenario
				repo.
					EXPECT().
					GetVideoByID(gomock.Any(), gomock.Any()).
					MaxTimes(0)

				_, err := videoService.GetByID(context.Background(), c.videoID)

				require.Error(t, err)
			} else {
				// Non-Error scenario
				repo.
					EXPECT().
					GetVideoByID(gomock.Any(), gomock.Eq(c.videoID)).
					Return(expectedVideo, nil).
					MaxTimes(1)

				res, err := videoService.GetByID(context.Background(), c.videoID)

				require.NoError(t, err)
				require.Equal(t, expectedVideo, res)
			}
		})
	}
}

// go test ./internal/app -run TestGetVideoByPublisher -v
func TestGetVideoByPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var emptyPublisherID = uuid.Nil
	var expectedRes = []domain.Video{{
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: testDesc,
		CreatedAt:   time.Now(),
	}}

	cases := []struct {
		name        string
		wantErr     bool
		publisherID uuid.UUID
		offset      int32
		limit       int32
		orderBy     string
		asc         string
	}{
		{"ok", false,
			testPublisherID, 5, 5, "createdAt", "t"},
		{"offset zero", false,
			testPublisherID, 0, 5, "createdAt", "t"},
		{"offset less zero", true,
			testPublisherID, -10, 5, "createdAt", "t"},
		{"limit less zero", true,
			testPublisherID, 0, -10, "createdAt", "t"},
		{"limit reached max value", true,
			testPublisherID, 0, 1 >> 8, "createdAt", "t"},
		{"invalid orderBy", true,
			testPublisherID, 0, 5, "crtdAt", "t"},
		{"invalid asc", true,
			testPublisherID, 0, 5, "createdAt", "wrong_asc"},
		{"nil publisher id", true,
			emptyPublisherID, 0, 5, "createdAt", "t"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				repo.
					EXPECT().
					GetPublisherVideos(gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
				_, err := videoService.GetByPublisher(
					context.Background(),
					c.publisherID,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.Error(t, err)
			} else {
				repo.
					EXPECT().
					GetPublisherVideos(gomock.Any(), gomock.Eq(c.publisherID), gomock.Any()).
					Return(expectedRes, nil).
					MaxTimes(1)
				res, err := videoService.GetByPublisher(
					context.Background(),
					c.publisherID,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.NoError(t, err)
				require.Exactly(t, expectedRes, res)
			}
		})
	}
}

func TestSearchVideoByPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var emptyPublisherID = uuid.Nil
	var expectedRes = []domain.Video{{
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: testDesc,
		CreatedAt:   time.Now(),
	}}

	cases := []struct {
		name        string
		wantErr     bool
		publisherID uuid.UUID
		query       string
		offset      int32
		limit       int32
		orderBy     string
		asc         string
	}{
		{"ok", false, testPublisherID, "search",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"ok with search surrounded spaces", false, testPublisherID, "   ok with spacing search    ",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"ok zero offset", false, testPublisherID, "search",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"offset less zero", true, testPublisherID, "search",
			-1 /*offset*/, 5 /*limit*/, "createdAt", "t"}, // Fix: Change wantErr to true

		{"limit less zero", true, testPublisherID, "search",
			0 /*offset*/, -1 /*limit*/, "createdAt", "t"}, // Fix: Change wantErr to true

		{"limit zero", true, testPublisherID, "search",
			0 /*offset*/, 0 /*limit*/, "createdAt", "t"},

		{"limit reached max value", true, testPublisherID, "search",
			0 /*offset*/, 0 /*limit*/, "createdAt", "t"},

		{"nil publisher id", true, emptyPublisherID, "search",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"incorrect search", true, testPublisherID, "se!arch",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				repo.
					EXPECT().
					SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
				_, err := videoService.SearchPublisher(
					context.Background(),
					c.publisherID,
					c.query,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.Error(t, err)
			} else {
				repo.
					EXPECT().
					SearchPublisher(gomock.Any(), gomock.Eq(c.publisherID), gomock.Any()).
					Return(expectedRes, nil).
					MaxTimes(1)
				res, err := videoService.SearchPublisher(
					context.Background(),
					c.publisherID,
					c.query,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.NoError(t, err)
				require.Exactly(t, expectedRes, res)
			}
		})
	}
}

func TestValidSearchGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var expectedRes = []domain.Video{{
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: testDesc,
		CreatedAt:   time.Now(),
	}}

	cases := []struct {
		name    string
		wantErr bool
		query   string
		offset  int32
		limit   int32
		orderBy string
		asc     string
	}{
		{"ok", false, "search",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"ok with search surrounded spaces", false, "   ok with spacing search    ",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},

		{"offset less zero", true, "search",
			-1 /*offset*/, 5 /*limit*/, "createdAt", "t"}, // Fix: Change wantErr to true

		{"limit less zero", true, "search",
			0 /*offset*/, -1 /*limit*/, "createdAt", "t"}, // Fix: Change wantErr to true

		{"limit zero", true, "search",
			0 /*offset*/, 0 /*limit*/, "createdAt", "t"},

		{"limit reached max value", true, "search",
			0 /*offset*/, 0 /*limit*/, "createdAt", "t"},

		{"incorrect search", true, "se!arch",
			0 /*offset*/, 5 /*limit*/, "createdAt", "t"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				repo.
					EXPECT().
					SearchGlobal(gomock.Any(), gomock.Any()).
					MaxTimes(0)
				_, err := videoService.SearchGlobal(
					context.Background(),
					c.query,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.Error(t, err)
			} else {
				repo.
					EXPECT().
					SearchGlobal(gomock.Any(), gomock.Any()).
					Return(expectedRes, nil).
					MaxTimes(1)
				res, err := videoService.SearchGlobal(
					context.Background(),
					c.query,
					c.orderBy,
					c.offset,
					c.limit,
					c.asc)
				require.NoError(t, err)
				require.Exactly(t, expectedRes, res)
			}
		})
	}
}
