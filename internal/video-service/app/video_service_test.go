package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"video-service/internal/app"
	"video-service/internal/domain"
	"video-service/internal/policy"
	"video-service/internal/ports"
	mock_ports "video-service/internal/ports/mock"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

var testTopic = "topic"
var testDesc = "desc"
var testPublisherID, _ = uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
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

				require.Equal(t, expectedVideo, res)
				require.NoError(t, err)
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
		name               string
		wantErr            bool
		publisherID        uuid.UUID
		pagination         ports.PageRequest
		expectedPagination ports.PageRequest
	}{
		{"ok", false,
			testPublisherID, getPageRequest(5, 5), getPageRequest(5, 5)},
		{"ok without offset", false,
			testPublisherID, getPageRequest(0, 5), getPageRequest(0, 5)},
		{"limit less zero pagination", false,
			testPublisherID, getPageRequest(0, -1), getPageRequest(0, policy.DefaultVideosLimitPerRequest)},
		{"offset less zero pagination", false,
			testPublisherID, getPageRequest(-1, 5), getPageRequest(0, 5)},
		{"limit zero pagination", false,
			testPublisherID, getPageRequest(5, -1), getPageRequest(5, policy.DefaultVideosLimitPerRequest)},
		{"nil publisher id", true,
			emptyPublisherID, getPageRequest(5, 5), getPageRequest(0, 0) /* error */},
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
					c.pagination.Offset,
					c.pagination.Limit)
				require.Error(t, err)
			} else {
				repo.
					EXPECT().
					GetPublisherVideos(gomock.Any(), gomock.Eq(c.publisherID), gomock.Eq(c.expectedPagination)).
					Return(expectedRes, nil).
					MaxTimes(1)
				res, err := videoService.GetByPublisher(
					context.Background(),
					c.publisherID,
					c.pagination.Offset,
					c.pagination.Limit)
				require.Exactly(t, expectedRes, res)
				require.NoError(t, err)
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
		search      ports.VideoSearch
	}{
		{"ok", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"ok with search surrounded spaces", false, testPublisherID,
			ports.VideoSearch{
				PageRequest: getPageRequest(5, 5),
				Query:       "   ok with spacing search    ",
			}},
		{"ok without offset", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search"}},
		{"limit less zero pagination", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search"}},
		{"offset less zero pagination", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search"}},
		{"limit zero pagination", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search"}},

		{"nil publisher id", true, emptyPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"incorrect search", true, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "se!arch"}},
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
					c.search.Query,
					c.search.Offset,
					c.search.Limit)
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
					c.search.Query,
					c.search.Offset,
					c.search.Limit)
				require.Exactly(t, expectedRes, res)
				require.NoError(t, err)
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
		search  ports.VideoSearch
	}{
		{"ok", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search global"}},
		{"ok with search surrounded spaces", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing    "}},
		{"ok without offset", false,
			ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search global"}},
		{"limit less zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search global"}},
		{"offset less zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search global"}},
		{"limit zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search global"}},
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
					c.search.Query,
					c.search.Offset,
					c.search.Limit)
				require.Error(t, err)
			} else {
				repo.
					EXPECT().
					SearchGlobal(gomock.Any(), gomock.Any()).
					Return(expectedRes, nil).
					MaxTimes(1)
				res, err := videoService.SearchGlobal(
					context.Background(),
					c.search.Query,
					c.search.Offset,
					c.search.Limit)
				require.NoError(t, err)
				require.Exactly(t, expectedRes, res)
			}
		})
	}
}

func getPageRequest(offset, limit int32) ports.PageRequest {
	return ports.PageRequest{Offset: offset, Limit: limit}
}
