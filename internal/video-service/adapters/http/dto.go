package httpadapter

import (
	"time"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/policy"

	"github.com/go-playground/validator/v10"
)

const (
	IDEmpty                 string = "ID_EMPTY"
	IDSizeExceeded          string = "ID_SIZE_EXCEEDED"
	TopicEmpty              string = "TOPIC_EMPTY"
	TopicSizeExceeded       string = "TOPIC_SIZE_EXCEEDED"
	DescriptionEmpty        string = "DESCRIPTION_EMPTY"
	DescriptionSizeExceeded string = "DESCRIPTION_SIZE_EXCEEDED"
)

type VideoResponseBody struct {
	ID          string    `json:"id"`
	PublisherID string    `json:"publisherID"`
	Topic       string    `json:"topic"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ValidationError struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (e ValidationError) Error() string {
	return e.ErrorMessage
}

// createVideoRequestBody represents the data required to create the video
type createVideoRequestBody struct {
	Topic       string `json:"topic" validate:"required,minTopic,maxTopic"`
	Description string `json:"description" validate:"required,maxDescription"`
}

func dtoVideo(v domain.Video) VideoResponseBody {
	return VideoResponseBody{
		ID:          v.ID.String(),
		PublisherID: v.PublisherID.String(),
		Topic:       v.Topic,
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

func dtoVideos(videos []domain.Video) []VideoResponseBody {
	res := make([]VideoResponseBody, len(videos))
	for i, v := range videos {
		res[i] = dtoVideo(v)
	}
	return res
}

// validate validates the request body fields.
// It checks for empty fields, length constraints, and returns an error.
func (r createVideoRequestBody) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("maxTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.MaxTopicLen
	})
	validate.RegisterValidation("minTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= policy.MinTopicLen
	})
	validate.RegisterValidation("maxDescription", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.MaxDescriptionLen
	})

	err := validate.Struct(r)
	if err != nil {
		switch err := err.(type) {
		case validator.ValidationErrors:
			return err[0]
		default:
			return err
		}
	}

	return nil
}
