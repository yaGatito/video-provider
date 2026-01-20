package domain

import (
	"errors"
	"video-service/internal/policy"

	"github.com/google/uuid"
)

func (v Video) Validate() error {
	if v.PublisherID == uuid.Nil {
		return errors.New("publisher id is empty")
	}
	if v.Topic == "" {
		return errors.New("topic is empty")
	}
	if !policy.GET_TEXTING_FORMAT_RE_128().MatchString(v.Topic) {
		return errors.New("invalid topic text format")
	}
	if v.Description != nil && *v.Description != "" && !policy.GET_LARGE_TEXT_FORMAT_RE_512().MatchString(*v.Description) {
		return errors.New("invalid description text format")
	}
	return nil
}
