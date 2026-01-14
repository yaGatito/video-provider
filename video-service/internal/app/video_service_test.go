package app

import (
	"context"
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

// go test ./internal/app -run TestGetdVideoByPublisher -v
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

	var pagination = ports.PageRequest{Offset: 5, Limit: 5}
	var expectedPagination = ports.PageRequest{Offset: 5, Limit: 5}

	var paginationWithoutOffset = ports.PageRequest{Limit: 5}
	var expectedPaginationWithoutOffset = ports.PageRequest{Offset: 0, Limit: 5}

	var negativeLimitPagination = ports.PageRequest{Limit: -1}
	var expectedNegativeLimitPagination = ports.PageRequest{Offset: 0, Limit: policy.MAX_VIDEOS_LIMIT_PER_REQUEST}

	var negativeOffsetPagination = ports.PageRequest{Offset: -1, Limit: 5}
	var expectedNegativeOffsetPagination = ports.PageRequest{Offset: 0, Limit: 5}

	var limitZeroPagination = ports.PageRequest{}
	var expectedLimitZeroPagination = ports.PageRequest{Offset: 0, Limit: policy.MAX_VIDEOS_LIMIT_PER_REQUEST}

	tests := []struct {
		name               string
		wantErr            bool
		publisherID        uuid.UUID
		pagination         ports.PageRequest
		expectedPagination ports.PageRequest
	}{
		{"ok", WANT_NO_ERROR, publisherID, pagination, expectedPagination},
		{"ok without offset", WANT_NO_ERROR, publisherID, paginationWithoutOffset, expectedPaginationWithoutOffset},
		{"limit less zero pagination", WANT_NO_ERROR, publisherID, negativeLimitPagination, expectedNegativeLimitPagination},
		{"offset less zero pagination", WANT_NO_ERROR, publisherID, negativeOffsetPagination, expectedNegativeOffsetPagination},
		{"limit zero pagination", WANT_NO_ERROR, publisherID, limitZeroPagination, expectedLimitZeroPagination},

		{"nil publisher id", WANT_ERROR, emptyPublisherID, pagination, ports.PageRequest{} /* any */},
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
