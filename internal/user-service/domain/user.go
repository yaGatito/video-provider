package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	LastName  string
	CreatedAt time.Time // 2025-08-14 00:37:00
	IsAdmin   bool      // internal-only; guarded by service logic
	Status    string    // "active", "disabled"
}
