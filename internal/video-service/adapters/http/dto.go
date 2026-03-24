package httpadp

import (
	"time"
	"video-provider/internal/video-service/domain"

	"github.com/go-playground/validator/v10"
)

type videoResponseBody struct {
	ID          string `json:"id"`
	PublisherID string `json:"publisherID"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type videosResponseBody struct {
	Videos []videoResponseBody `json:"videos"`
}

type serviceErrorResponse struct {
	Message string `json:"msg"`
}

// createVideoRequestBody represents the data required to create the video
type createVideoRequestBody struct {
	Topic       string `json:"topic" validate:"required,minTopic,maxTopic"`
	Description string `json:"description" validate:"required,maxDescription"`
}

func (r createVideoRequestBody) validate(v *validator.Validate) error {
	err := v.Struct(r)
	if err != nil {
		return err
	}

	return nil
}

func dtoVideo(v domain.Video) videoResponseBody {
	return videoResponseBody{
		ID:          v.ID.String(),
		PublisherID: v.PublisherID.String(),
		Topic:       v.Topic,
		Description: v.Description,
		CreatedAt:   v.CreatedAt.Format(time.DateTime),
	}
}
