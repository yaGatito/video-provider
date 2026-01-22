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

var testTopic = "topic"
var testDesc = "desc"
var testPublisherID, _ = uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
var testVideo = domain.Video{
	PublisherID: testPublisherID,
	Topic:       testTopic,
	Description: &testDesc,
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

	err := videoService.Create(context.Background(), testVideo)

	if err != nil {
		t.Fatalf("name: ok; input data: %v, res error %e\n", testVideo, err)
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
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: &testDesc,
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		videoID uuid.UUID
		wantErr bool
	}{
		{"ok", testVideoID, false},
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
				GetVideoByID(gomock.Any(), gomock.Eq(tt.videoID)).
				Return(expectedVideo, nil).
				MaxTimes(1)
		}
		res, err := videoService.GetByID(context.Background(), tt.videoID)
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
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: &testDesc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
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
			testPublisherID, getPageRequest(0, -1), getPageRequest(0, policy.MaxVideosLimitPerRequest)},
		{"offset less zero pagination", false,
			testPublisherID, getPageRequest(-1, 5), getPageRequest(0, 5)},
		{"limit zero pagination", false,
			testPublisherID, getPageRequest(5, -1), getPageRequest(5, policy.MaxVideosLimitPerRequest)},
		{"nil publisher id", true,
			emptyPublisherID, getPageRequest(5, 5), getPageRequest(0, 0) /* error */},
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
				GetPublisherVideos(gomock.Any(), gomock.Eq(tt.publisherID), gomock.Eq(tt.expectedPagination)).
				Return(expectedRes, nil).
				MaxTimes(1)
		}

		res, err := videoService.GetByPublisher(
			context.Background(),
			tt.publisherID,
			tt.pagination.Offset,
			tt.pagination.Limit)

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
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: &testDesc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name        string
		wantErr     bool
		publisherID uuid.UUID
		search      ports.VideoSearch
	}{
		{"ok", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"ok with search surrounded spaces", false, testPublisherID,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
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

		res, err := videoService.SearchPublisher(
			context.Background(),
			tt.publisherID,
			tt.search.Query,
			tt.search.Offset,
			tt.search.Limit)

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
		ID:          testVideoID,
		PublisherID: testPublisherID,
		Topic:       testTopic,
		Description: &testDesc,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name    string
		wantErr bool
		search  ports.VideoSearch
	}{
		{"ok", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search global"}},
		{"ok with search surrounded spaces", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
		{"ok without offset", false,
			ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search global"}},
		{"limit less zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search global"}},
		{"offset less zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search global"}},
		{"limit zero pagination", false,
			ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search global"}},
	}

	for _, tt := range tests {
		repo.
			EXPECT().
			SearchGlobal(gomock.Any(), gomock.Any()).
			Return(expectedRes, nil).
			MaxTimes(1)

		res, err := videoService.SearchGlobal(
			context.Background(),
			tt.search.Query,
			tt.search.Offset,
			tt.search.Limit)

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

		_, err := videoService.SearchGlobal(
			context.Background(),
			tt.search.Query,
			tt.search.Offset,
			tt.search.Limit)

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
