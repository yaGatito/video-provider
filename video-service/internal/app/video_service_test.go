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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

var test_topic = "test_topic"
var test_desc = "desc"
var test_publisherID, _ = uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
var test_video = domain.Video{
	PublisherID: test_publisherID,
	Topic:       "topic",
	Description: &test_desc,
}
var test_videoID, _ = uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")

func TestCreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	repo.
		EXPECT().
		CreateVideo(gomock.Any(), gomock.Eq(test_video)).
		MaxTimes(1)

	err := videoService.Create(context.Background(), test_video)

	if err != nil {
		t.Fatalf("name: ok; input data: %v, res error %e\n", test_video, err)
		return
	}
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

	err := videoService.Create(context.Background(), invalidVideo)

	if err == nil {
		t.Fatalf("name: invalid video; input data: %v, res error %e\n", invalidVideo, err)
		return
	}
}

func TestGetdVideoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	expectedVideo := domain.Video{
		ID:          test_videoID,
		PublisherID: test_publisherID,
		Topic:       test_topic,
		Description: &test_desc,
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name         string
		test_videoID uuid.UUID
		wantErr      bool
	}{
		{"ok", test_videoID, false},
		{"no video id", uuid.Nil, true},
	}

	for _, tt := range tests {
		if tt.wantErr == true {
			repo.
				EXPECT().
				GetVideoByID(gomock.Any(), gomock.Any()).
				MaxTimes(0)
		} else {
			repo.
				EXPECT().
				GetVideoByID(gomock.Any(), gomock.Eq(tt.test_videoID)).
				Return(expectedVideo, nil).
				MaxTimes(1)
		}
		res, err := videoService.GetByID(context.Background(), tt.test_videoID)
		isEmptyError := err == nil
		if tt.wantErr == isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
		if isEmptyError && !gomock.Eq(expectedVideo).Matches(res) {
			t.Fatalf("res:%v, expected:%v", res, expectedVideo)
			return
		}
	}
}

// go test ./internal/app -run TestGetVideoByPublisher -v
func TestGetVideoByPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var emptyPublisherID = uuid.Nil
	var expectedRes = []domain.Video{{
		ID:          test_videoID,
		PublisherID: test_publisherID,
		Topic:       test_topic,
		Description: &test_desc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name               string
		wantErr            bool
		test_publisherID   uuid.UUID
		pagination         ports.PageRequest
		expectedPagination ports.PageRequest
	}{
		{"ok", false, test_publisherID, getPageRequest(5, 5), getPageRequest(5, 5)},
		{"ok without offset", false, test_publisherID, getPageRequest(0, 5), getPageRequest(0, 5)},
		{"limit less zero pagination", false, test_publisherID, getPageRequest(0, -1), getPageRequest(0, policy.MAX_VIDEOS_LIMIT_PER_REQUEST)},
		{"offset less zero pagination", false, test_publisherID, getPageRequest(-1, 5), getPageRequest(0, 5)},
		{"limit zero pagination", false, test_publisherID, getPageRequest(5, -1), getPageRequest(5, policy.MAX_VIDEOS_LIMIT_PER_REQUEST)},

		{"nil publisher id", true, emptyPublisherID, getPageRequest(5, 5), getPageRequest(0, 0) /* error */},
	}

	for _, tt := range tests {
		if tt.wantErr {
			repo.
				EXPECT().
				GetPublisherVideos(gomock.Any(), gomock.Any(), gomock.Any()).
				MaxTimes(0)
		} else {
			repo.
				EXPECT().
				GetPublisherVideos(gomock.Any(), gomock.Eq(tt.test_publisherID), gomock.Eq(tt.expectedPagination)).
				Return(expectedRes, nil).
				MaxTimes(1)
		}

		res, err := videoService.GetByPublisher(context.Background(), tt.test_publisherID, tt.pagination.Offset, tt.pagination.Limit)
		isEmptyError := err == nil
		if tt.wantErr == isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
		if isEmptyError && !gomock.Eq(expectedRes).Matches(res) {
			t.Fatalf("res:%v, expected:%v", res, expectedRes)
			return
		}
	}
}

func TestSearchVideoByPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var emptyPublisherID = uuid.Nil
	var expectedRes = []domain.Video{{
		ID:          test_videoID,
		PublisherID: test_publisherID,
		Topic:       test_topic,
		Description: &test_desc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name        string
		wantErr     bool
		publisherID uuid.UUID
		search      ports.VideoSearch
	}{
		{"ok", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"ok with search surrounded spaces", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
		{"ok without offset", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search"}},
		{"limit less zero pagination", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search"}},
		{"offset less zero pagination", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search"}},
		{"limit zero pagination", false, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search"}},

		{"nil publisher id", true, emptyPublisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"incorrect search", true, test_publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "se!arch"}},
	}

	for _, tt := range tests {
		if tt.wantErr {
			repo.
				EXPECT().
				SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any()).
				MaxTimes(0)
		} else {
			repo.
				EXPECT().
				SearchPublisher(gomock.Any(), gomock.Eq(tt.publisherID), gomock.Any()).
				Return(expectedRes, nil).
				MaxTimes(1)
		}

		res, err := videoService.SearchPublisher(context.Background(), tt.publisherID, tt.search.Query, tt.search.Offset, tt.search.Limit)
		isEmptyError := err == nil
		if tt.wantErr == isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
		if isEmptyError && !gomock.Eq(expectedRes).Matches(res) {
			t.Fatalf("res:%v, expected:%v", res, expectedRes)
			return
		}
	}
}

func TestValidSearchGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	var expectedRes = []domain.Video{{
		ID:          test_videoID,
		PublisherID: test_publisherID,
		Topic:       test_topic,
		Description: &test_desc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name    string
		wantErr bool
		search  ports.VideoSearch
	}{
		{"ok", false, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search global"}},
		{"ok with search surrounded spaces", false, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
		{"ok without offset", false, ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search global"}},
		{"limit less zero pagination", false, ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search global"}},
		{"offset less zero pagination", false, ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search global"}},
		{"limit zero pagination", false, ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search global"}},
	}

	for _, tt := range tests {
		repo.
			EXPECT().
			SearchGlobal(gomock.Any(), gomock.Any()).
			Return(expectedRes, nil).
			MaxTimes(1)

		res, err := videoService.SearchGlobal(context.Background(), tt.search.Query, tt.search.Offset, tt.search.Limit)
		isEmptyError := err == nil
		if tt.wantErr == isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
		if isEmptyError && !gomock.Eq(expectedRes).Matches(res) {
			t.Fatalf("res:%v, expected:%v", res, expectedRes)
			return
		}
	}
}

func TestIncorrectSearchGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := app.NewVideoInteractor(repo)

	tests := []struct {
		name    string
		wantErr bool
		search  ports.VideoSearch
	}{}

	var symbols = "@#$%^&*()+="
	var format = "se%carch global"
	for i, c := range symbols {
		newtt := struct {
			name    string
			wantErr bool
			search  ports.VideoSearch
		}{
			name:    fmt.Sprintf("incorrect search: %d; symbol: %c", i+1, c),
			wantErr: true,
			search:  ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: fmt.Sprintf(format, c)},
		}

		tests = append(tests, newtt)
	}

	for _, tt := range tests {
		repo.
			EXPECT().
			SearchGlobal(gomock.Any(), gomock.Any()).
			MaxTimes(0)

		_, err := videoService.SearchGlobal(context.Background(), tt.search.Query, tt.search.Offset, tt.search.Limit)
		isEmptyError := err == nil
		if isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
	}
}

func getPageRequest(offset, limit int32) ports.PageRequest {
	return ports.PageRequest{Offset: offset, Limit: limit}
}
