package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

type UUID = uuid.UUID

type Status string

type Video struct {
	ID          UUID
	PublisherID UUID
	Topic       string
	Description *string
	CreatedAt   time.Time
	Status      Status
}

func (v Video) Validate() error {
	if v.PublisherID == uuid.Nil {
		return fmt.Errorf("publisher id is empty")
	}
	if v.Topic == "" {
		return fmt.Errorf("topic is empty")
	}
	return nil
}
