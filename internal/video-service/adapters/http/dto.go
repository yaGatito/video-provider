package httpadp

import (
	"time"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/policy"

	"github.com/go-playground/validator/v10"
)

type videoResponseBody struct {
	ID          string `json:"id"`
	PublisherID string `json:"publisherID"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

// createVideoRequestBody represents the data required to create the video
type createVideoRequestBody struct {
	Topic       string `json:"topic" validate:"required,minTopic,maxTopic"`
	Description string `json:"description" validate:"required,maxDescription"`
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

// validate validates the request body fields.
// It checks for empty fields, length constraints, and returns an error.
func (r createVideoRequestBody) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("maxTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.TopicMaxLen
	})
	validate.RegisterValidation("minTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= policy.TopicMinLen
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

	if _, ok := err.(validator.ValidationErrors); ok {

	}

	return nil
}
