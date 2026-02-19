package domain

import (
	"errors"
	"video-provider/internal/video-service/policy"

	"github.com/google/uuid"
)

func (v Video) Validate() error {
	if v.PublisherID == uuid.Nil {
		return errors.New("publisher id is empty")
	}
	if v.Topic == "" {
		return errors.New("topic is empty")
	}
	if !policy.GetTextingFormateRE128().MatchString(v.Topic) {
		return errors.New("invalid topic text format")
	}
	if v.Description != "" &&
		!policy.GetLargeTextFormatRE512().MatchString(v.Description) {
		return errors.New("invalid description text format")
	}
	return nil
}
