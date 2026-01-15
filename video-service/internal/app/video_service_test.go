package app

import (
	"context"
	"fmt"
	"testing"
	"time"
	"video-service/internal/domain"
	"video-service/internal/policy"
	"video-service/internal/ports"
	mock_ports "video-service/internal/ports/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

const WANT_NO_ERROR = false
const WANT_ERROR = true

var TEST_DATA_DESCRIPTION string = "TEST"

func TestCreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := NewVideoInteractor(repo)

	publisherID, _ := uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")

	tests := []struct {
		name    string
		wantErr bool
		video   domain.Video
	}{
		{"ok", WANT_NO_ERROR, domain.Video{
			PublisherID: publisherID,
			Topic:       "topic",
			Description: &TEST_DATA_DESCRIPTION},
		},
		{"ok - no desc", WANT_NO_ERROR, domain.Video{
			PublisherID: publisherID,
			Topic:       "topic"},
		},
		{"no topic", WANT_ERROR, domain.Video{
			PublisherID: publisherID},
		},
		{"no publisher id", WANT_ERROR, domain.Video{
			Topic: "topic"},
		},
	}

	for _, tt := range tests {
		if tt.wantErr {
			repo.
				EXPECT().
				CreateVideo(gomock.Any(), gomock.Any()).
				MaxTimes(0)
		} else {
			repo.
				EXPECT().
				CreateVideo(gomock.Any(), gomock.Eq(tt.video)).
				MaxTimes(1)
		}

		err := videoService.Create(context.Background(), tt.video)

		if (tt.wantErr && err == nil) || (!tt.wantErr && err != nil) {
			t.Fatalf("tt.wantErr %t; res error %e\n", tt.wantErr, err)
			return
		}
	}
}

func TestGetdVideoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := NewVideoInteractor(repo)

	videoID, _ := uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")
	publisherID, _ := uuid.Parse("464f522f-0106-d9fa-8d98-356bad96a56b")

	expectedVideo := domain.Video{
		ID:          videoID,
		PublisherID: publisherID,
		Topic:       "topic",
		Description: &TEST_DATA_DESCRIPTION,
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		videoID uuid.UUID
		wantErr bool
	}{
		{"ok", videoID, WANT_NO_ERROR},
		{"no video id", uuid.Nil, WANT_ERROR},
	}

	for _, tt := range tests {
		if tt.wantErr == WANT_ERROR {
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
	videoService := NewVideoInteractor(repo)

	var publisherID, _ = uuid.Parse("464f522f-0106-d9fa-8d98-356bad96a56b")
	var emptyPublisherID = uuid.Nil
	var videoID, _ = uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")
	var expectedRes = []domain.Video{{
		ID:          videoID,
		PublisherID: publisherID,
		Topic:       "topic",
		Description: &TEST_DATA_DESCRIPTION,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name               string
		wantErr            bool
		publisherID        uuid.UUID
		pagination         ports.PageRequest
		expectedPagination ports.PageRequest
	}{
		{"ok", WANT_NO_ERROR, publisherID, getPageRequest(5, 5), getPageRequest(5, 5)},
		{"ok without offset", WANT_NO_ERROR, publisherID, getPageRequest(0, 5), getPageRequest(0, 5)},
		{"limit less zero pagination", WANT_NO_ERROR, publisherID, getPageRequest(0, -1), getPageRequest(0, policy.MAX_VIDEOS_LIMIT_PER_REQUEST)},
		{"offset less zero pagination", WANT_NO_ERROR, publisherID, getPageRequest(-1, 5), getPageRequest(0, 5)},
		{"limit zero pagination", WANT_NO_ERROR, publisherID, getPageRequest(5, -1), getPageRequest(5, policy.MAX_VIDEOS_LIMIT_PER_REQUEST)},

		{"nil publisher id", WANT_ERROR, emptyPublisherID, getPageRequest(5, 5), getPageRequest(0, 0) /* error */},
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

		res, err := videoService.GetByPublisher(context.Background(), tt.publisherID, tt.pagination.Offset, tt.pagination.Limit)
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
	videoService := NewVideoInteractor(repo)

	var publisherID, _ = uuid.Parse("464f522f-0106-d9fa-8d98-356bad96a56b")
	var emptyPublisherID = uuid.Nil

	var videoID, _ = uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")
	var expectedRes = []domain.Video{{
		ID:          videoID,
		PublisherID: publisherID,
		Topic:       "topic",
		Description: &TEST_DATA_DESCRIPTION,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name        string
		wantErr     bool
		publisherID uuid.UUID
		search      ports.VideoSearch
	}{
		{"ok", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"ok with search surrounded spaces", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
		{"ok without offset", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search"}},
		{"limit less zero pagination", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search"}},
		{"offset less zero pagination", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search"}},
		{"limit zero pagination", WANT_NO_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search"}},

		{"nil publisher id", WANT_ERROR, emptyPublisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search"}},
		{"incorrect search", WANT_ERROR, publisherID, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "se!arch"}},
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
	videoService := NewVideoInteractor(repo)

	var publisherID, _ = uuid.Parse("464f522f-0106-d9fa-8d98-356bad96a56b")
	var videoID, _ = uuid.Parse("d9fa522f-0016-464f-8d68-356ba1d6ad7d")
	var expectedRes = []domain.Video{{
		ID:          videoID,
		PublisherID: publisherID,
		Topic:       "topic",
		Description: &TEST_DATA_DESCRIPTION,
		CreatedAt:   time.Now(),
	}}

	tests := []struct {
		name    string
		wantErr bool
		search  ports.VideoSearch
	}{
		{"ok", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "search global"}},
		{"ok with search surrounded spaces", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(5, 5), Query: "   ok with spacing search    "}},
		{"ok without offset", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(0, 5), Query: "search global"}},
		{"limit less zero pagination", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(0, -1), Query: "search global"}},
		{"offset less zero pagination", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(-1, 5), Query: "search global"}},
		{"limit zero pagination", WANT_NO_ERROR, ports.VideoSearch{PageRequest: getPageRequest(5, -1), Query: "search global"}},
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

func TestInvalidSearchGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_ports.NewMockVideoRepository(ctrl)
	videoService := NewVideoInteractor(repo)

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
			wantErr: WANT_ERROR,
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
