package app

import (
	"fmt"
	"regexp"
	"strings"
	"video-service/internal/policy"
)

var queryRe = regexp.MustCompile(`^[\p{L}\p{N}\s_-]{1,100}$`)

// ValidateSearchQuery Returns trimmed version of query.
func ValidateSearchQuery(query string) (string, error) {
	qBytes := []byte(strings.TrimSpace(query))

	if len(qBytes) == 0 {
		return "", fmt.Errorf("query len is zero")
	}
	if len(qBytes) > policy.MAX_SEARCH_BYTES_SIZE {
		return "", fmt.Errorf("query len more then limit %d bytes", policy.MAX_SEARCH_BYTES_SIZE)
	}
	if !queryRe.MatchString(string(qBytes)) {
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
		limit = policy.MAX_VIDEOS_LIMIT_PER_REQUEST
	}
	if limit > policy.MAX_VIDEOS_LIMIT_PER_REQUEST {
		limit = policy.MAX_VIDEOS_LIMIT_PER_REQUEST
	}
	return offset, limit
}
