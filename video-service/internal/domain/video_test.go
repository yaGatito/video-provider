package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestVideoValidate(t *testing.T) {
	publisherID, _ := uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	tests := []struct {
		name        string
		publisherID UUID
		topic       string
		description string
		wantErr     bool
	}{
		{"ok", publisherID, "topic", "desc", false},
		{"empty desc", publisherID, "topic", "", false},
		{"empty topic", publisherID, "", "desc", true},
		{"empty id", uuid.Nil, "topic", "desc", true},
	}

	for _, tt := range tests {
		v := Video{
			PublisherID: tt.publisherID,
			Topic:       tt.topic,
			Description: &tt.description,
		}
		err := v.Validate()
		if (err != nil) != tt.wantErr {
			t.Fatalf("%s: unexpected result", tt.name)
		}
	}
}

