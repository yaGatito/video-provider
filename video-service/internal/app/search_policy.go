package app

import (
	"fmt"
	"strings"
	"video-service/internal/policy"
)

// ValidateSearchQuery Returns trimmed version of query.
func ValidateSearchQuery(query string) (string, error) {
	qBytes := []byte(strings.TrimSpace(query))

	if len(qBytes) == 0 {
		return "", fmt.Errorf("query len is zero")
	}
	if len(qBytes) > policy.MaxSearchBytesSize {
		return "", fmt.Errorf("query len more then limit %d bytes", policy.MaxSearchBytesSize)
	}
	if len(qBytes) < policy.MinSearchBytesSize {
		return "", fmt.Errorf("query len less then limit %d bytes", policy.MinSearchBytesSize)
	}
	if !policy.GetWordsFormatRE128().MatchString(string(qBytes)) {
		return "", fmt.Errorf("query string contains prohibited characters")
	}

	return string(qBytes), nil
}

// ValidatePagination Returns provided or default values for offset, limit.
func ValidatePagination(offset, limit int32) (int32, int32) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = policy.MaxVideosLimitPerRequest
	}
	if limit > policy.MaxVideosLimitPerRequest {
		limit = policy.MaxVideosLimitPerRequest
	}
	return offset, limit
}
