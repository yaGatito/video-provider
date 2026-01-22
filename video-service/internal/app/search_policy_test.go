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
		{"ok", 0, 5, 0, 5},
		{"zero limit", 0, 0, 0, policy.MaxVideosLimitPerRequest},
		{"negative limit", 5, -1, 5, policy.MaxVideosLimitPerRequest},
		{"negative offset", -1, 5, 0, 5},
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
		{"ok", false, "search", "search"},
		{"ok beggining spaces", false, "search", "      search"},
		{"ok ending spaces", false, "search", "search      "},
		{"ok surrounding spaces", false, "search", "      search      "},
		{"too short search query 1 len", true, "s", ""},
		{"too short search query 2 len", true, "s1", ""},
		{"too long search query", true, "LOOOOOOOOOOOOOOOOOOOOOOOOOOO" +
			"OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG SEARCH QUERY", ""},
		{"incorrect search query", true, "SE%#ARCH$@", ""},
		{"incorrect search query 2", true, "S!ARCH", ""},
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
