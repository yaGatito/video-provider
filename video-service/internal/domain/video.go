package domain

import (
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
	Description string
	CreatedAt   time.Time
	Status      Status
}
