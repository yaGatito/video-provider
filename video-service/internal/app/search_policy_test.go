package app

import (
	"testing"
	"video-service/internal/policy"

	"github.com/golang/mock/gomock"
)

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name           string
		offset         int32
		limit          int32
		expectedOffset int32
		expectedLimit  int32
	}{
		{"ok", 0, 1, 0, 1},
		{"zero limit", 0, 0, 0, policy.MAX_VIDEOS_LIMIT_PER_REQUEST},
		{"negative limit", 1, -1, 1, policy.MAX_VIDEOS_LIMIT_PER_REQUEST},
	}

	for _, tt := range tests {
		offset, limit := ValidatePagination(tt.offset, tt.limit)
		if offset != tt.expectedOffset {
			t.Fatalf("offset result %d and expected value %d not equals", offset, tt.expectedOffset)
		}
		if limit != tt.expectedLimit {
			t.Fatalf("limit result %d and expected value %d not equals", limit, tt.expectedLimit)
		}
	}
}

func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name           string
		wantErr        bool
		outputExpected string
		query          string
	}{
		{"ok", false, "search", "     search"},
		{"ok2", false, "search", "search      "},
		{"ok2", false, "search", "   search     "},
		{"ok3", false, "search", "search"},
		{"wrong search query", true, "LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG SEARCH QUERY", ""},
		{"wrong search query", true, "SE%#ARCH$@", ""},
		{"wrong search query", true, "S!ARCH", ""},
	}

	for _, tt := range tests {
		searchRes, err := ValidateSearchQuery(tt.query)
		isEmptyError := err == nil
		if tt.wantErr == isEmptyError {
			t.Fatalf("want error:%t, error:%s", tt.wantErr, err)
			return
		}
		if isEmptyError && !gomock.Eq(searchRes).Matches(tt.outputExpected) {
			t.Fatalf("res:%v, expected:%v", searchRes, tt.outputExpected)
			return
		}
	}
}
