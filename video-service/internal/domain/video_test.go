package domain_test

import (
	"testing"
	"video-service/internal/domain"

	"github.com/google/uuid"
)

func TestVideoValidate(t *testing.T) {
	publisherID, _ := uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	desc := "desc"
	specSymbolsDesc := "desc....,,;;;11!!"
	wrongSpecSymbolsDesc := "desc<><>//>\\"

	tests := []struct {
		name    string
		wantErr bool
		video   domain.Video
	}{
		{"ok", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "topic",
			Description: desc},
		},
		{"ok - no desc", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "topic"},
		},
		{"topic and desc with symbols", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "top  1111i  c!?1  11",
			Description: specSymbolsDesc},
		},
		{"2 symbol topic", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "aa",
			Description: desc},
		},
		{"2 symbol and 1+ spec symbols topic", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "aa!",
			Description: desc},
		},
		{"words", false, domain.Video{
			PublisherID: publisherID,
			Topic:       "    Lorem ipsum dolor; lorem    ipsum dolor !!! LOREM ISPUM DOLOR.   ",
			Description: desc},
		},
		{"2 spec symbol topic", true, domain.Video{
			PublisherID: publisherID,
			Topic:       "!!",
			Description: desc},
		},
		{"1 spec symbol topic", true, domain.Video{
			PublisherID: publisherID,
			Topic:       "!",
			Description: desc},
		},
		{"1 symbol topic", true, domain.Video{
			PublisherID: publisherID,
			Topic:       "a",
			Description: desc},
		},
		{"non-text topic", true, domain.Video{
			PublisherID: publisherID,
			Topic:       "top1111ic@$##$^$",
			Description: desc},
		},
		{"non-text description", true, domain.Video{
			PublisherID: publisherID,
			Topic:       "topic",
			Description: wrongSpecSymbolsDesc},
		},
		{"no topic", true, domain.Video{
			PublisherID: publisherID},
		},
		{"no publisher id", true, domain.Video{
			Topic: "topic"},
		},
	}

	for _, tt := range tests {
		err := tt.video.Validate()
		if tt.wantErr == (err == nil) {
			t.Fatalf(
				"name: %s; inputData%v; want error: %t; got error: %e",
				tt.name,
				tt.video,
				tt.wantErr,
				err,
			)
		}
	}
}
