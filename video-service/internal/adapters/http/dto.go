package httpadapter

import (
	"fmt"
	"time"
	"video-service/internal/policy"
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
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

// validate validates the request body fields.
// It checks for empty fields, length constraints, and returns an error.
func (r createVideoRequestBody) validate() error {
	topicSize := len([]byte(r.Topic))
	if topicSize == 0 {
		return ValidationError{
			ErrorCode:    TopicEmpty,
			ErrorMessage: "topic is empty",
		}
	}
	if topicSize > policy.MaxTopicBytesSize {
		return ValidationError{
			ErrorCode:    TopicSizeExceeded,
			ErrorMessage: fmt.Sprintf("topic size is more then %d", policy.MaxTopicBytesSize),
		}
	}

	descSize := len([]byte(r.Description))
	if descSize == 0 {
		return ValidationError{
			ErrorCode:    DescriptionEmpty,
			ErrorMessage: "description is empty",
		}
	}
	if descSize > policy.MaxDescriptionBytesSize {
		return ValidationError{
			ErrorCode:    DescriptionSizeExceeded,
			ErrorMessage: fmt.Sprintf("description size is more then %d", policy.MaxDescriptionBytesSize),
		}
	}

	return nil
}
