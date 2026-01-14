package crypto

import (
	"video-service/internal/ports"

	"github.com/google/uuid"
)

type UUIDGen struct {
}

var _ ports.IDGen = (*UUIDGen)(nil)

func New() UUIDGen {
	return UUIDGen{}
}

func (g UUIDGen) NewID() uuid.UUID {
	return uuid.New()
}

func (g UUIDGen) Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// 0d53e781-be85-4b67-a7d9-954c1271211c
// 0afb8a7a-336e-4a7a-ba66-676558d4833c
