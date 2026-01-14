package httpadp

import (
	"fmt"
	"time"
	"video-service/internal/policy"
)

const (
	ID_EMPTY                  string = "ID_EMPTY"
	ID_SIZE_EXCEEDED          string = "ID_SIZE_EXCEEDED"
	TOPIC_EMPTY               string = "TOPIC_EMPTY"
	TOPIC_SIZE_EXCEEDED       string = "TOPIC_SIZE_EXCEEDED"
	DESCRIPTION_EMPTY         string = "DESCRIPTION_EMPTY"
	DESCRIPTION_SIZE_EXCEEDED string = "DESCRIPTION_SIZE_EXCEEDED"
)

type VideoResponseBody struct {
	ID          string    `json:"id"`
	PublisherID string    `json:"publisherId"`
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
			ErrorCode: TOPIC_EMPTY, ErrorMessage: "topic is empty",
		}
	}
	if topicSize > policy.MAX_TOPIC_BYTES_SIZE {
		return ValidationError{
			ErrorCode: TOPIC_SIZE_EXCEEDED, ErrorMessage: fmt.Sprintf("topic size is more then %d", policy.MAX_TOPIC_BYTES_SIZE),
		}
	}

	descSize := len([]byte(r.Description))
	if descSize == 0 {
		return ValidationError{
			ErrorCode: DESCRIPTION_EMPTY, ErrorMessage: "description is empty",
		}
	}
	if descSize > policy.MAX_DESCRIPTION_BYTES_SIZE {
		return ValidationError{
			ErrorCode: DESCRIPTION_SIZE_EXCEEDED, ErrorMessage: fmt.Sprintf("description size is more then %d", policy.MAX_DESCRIPTION_BYTES_SIZE),
		}
	}

	return nil
}
