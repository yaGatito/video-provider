package ports

import (
	"video-provider/internal/video-service/domain"
)

type IDGen interface {
	NewID() domain.UUID
	Parse(s string) (domain.UUID, error)
}
