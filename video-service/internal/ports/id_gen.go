package ports

import (
	"video-service/internal/domain"
)

type IDGen interface {
	NewID() domain.UUID
	Parse(s string) (domain.UUID, error)
}
