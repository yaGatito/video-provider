package httpadp

import (
	"time"
	"video-service/domain"
)

type VideoResponseBody struct {
	ID          string `json:"id"`
	PublisherID string `json:"publisherID"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type VideosResponseBody struct {
	Videos []VideoResponseBody `json:"videos"`
}

type serviceErrorResponse struct {
	Message string `json:"msg"`
}

// createVideoRequestBody represents the data required to create the video
type createVideoRequestBody struct {
	Topic       string `json:"topic"       validate:"required,minTopic,maxTopic"`
	Description string `json:"description" validate:"required,maxDescription"`
}

// TODO: make it private
func DtoVideo(v domain.Video) VideoResponseBody {
	return VideoResponseBody{
		ID:          v.ID.String(),
		PublisherID: v.PublisherID.String(),
		Topic:       v.Topic,
		Description: v.Description,
		CreatedAt:   v.CreatedAt.Format(time.DateTime),
	}
}
