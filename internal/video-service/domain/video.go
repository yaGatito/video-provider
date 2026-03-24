package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"

	OrderByDate string = "date"
	AscOrder    string = "t"
	DescOrder   string = "f"
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

type VideoPageParams struct {
	OrderBy string
	Offset  int32
	Limit   int32
	Asc     string
}
