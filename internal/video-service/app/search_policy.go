package app

import (
	"fmt"
	"strings"
	"video-provider/internal/video-service/policy"
	"video-provider/internal/video-service/ports"
)

// ValidateSearchQuery trims the input query and validates it according to specific rules. It returns the trimmed query if valid, or an error if the query fails any of the validation checks.
//
// Parameters:
// - query: The input string to be validated and trimmed.
//
// Returns:
// - string: The trimmed version of the query if valid.
// - error: An error if the query fails validation, with details on the reason for failure.
//
// Validation Rules:
// 1. The query must not be empty after trimming. If it is, an error is returned indicating that the query length is zero.
// 2. The query must not exceed the maximum allowed size, defined by policy.MaxSearchBytesSize. If it does, an error is returned indicating that the query length exceeds the limit.
// 3. The query must not be shorter than the minimum allowed size, defined by policy.MinSearchBytesSize. If it does, an error is returned indicating that the query length is below the limit.
// 4. The query must only contain allowed characters, as defined by the regular expression pattern in policy.GetWordsFormatRE128(). If it contains prohibited characters, an error is returned indicating this violation.
func ValidateSearchQuery(query string) (string, error) {
	qBytes := []byte(strings.TrimSpace(query))

	if len(qBytes) == 0 {
		return "", fmt.Errorf("query len is zero")
	}
	if len(qBytes) > policy.MaxSearchBytesSize {
		return "", fmt.Errorf("query len more than limit %d bytes", policy.MaxSearchBytesSize)
	}
	if len(qBytes) < policy.MinSearchBytesSize {
		return "", fmt.Errorf("query len less than limit %d bytes", policy.MinSearchBytesSize)
	}
	if !policy.GetWordsFormatRE128().MatchString(string(qBytes)) {
		return "", fmt.Errorf("query string contains prohibited characters")
	}

	return string(qBytes), nil
}

// ValidateLimit validates the provided limit value to ensure it meets the required constraints for video retrieval requests. It returns the validated limit if it is within the allowed range, or an error if it fails any of the validation checks.
//
// Parameters:
// - limit: The integer value representing the number of videos to retrieve per request.
//
// Returns:
// - int32: The validated limit value if it is within the allowed range.
// - error: An error if the limit is invalid, with details on the reason for failure.
//
// Validation Rules:
// 1. The limit must be greater than the default value defined by policy.DefaultVideosLimitPerRequest. If it is less than or equal to this value, an error is returned indicating that the limit is zero or less.
// 2. The limit must not exceed the maximum allowed value defined by policy.MaxVideosLimitPerRequest. If it does, an error is returned indicating that the limit has reached the maximum allowed value.
func ValidateLimit(limit int32) (int32, error) {
	if limit <= policy.DefaultVideosLimitPerRequest {
		return 0, fmt.Errorf("limit is zero or less")
	}
	if limit > policy.MaxVideosLimitPerRequest {
		return 0, fmt.Errorf("limit reached maximum allowed value")
	}
	return limit, nil
}

// ValidateOffset validates the provided offset value to ensure it meets the required constraints for video retrieval requests. It returns the validated offset if it is within the allowed range, or an error if it fails any of the validation checks.
//
// Parameters:
// - offset: The integer value representing the starting position for retrieving videos.
//
// Returns:
// - int32: The validated offset value if it is within the allowed range.
// - error: An error if the offset is invalid, with details on the reason for failure.
//
// Validation Rules:
// 1. The offset must not be less than zero. If it is, an error is returned indicating that the offset is zero or less.
// 2. The offset must not exceed the maximum allowed value defined by policy.MaxInt32. If it does, an error is returned indicating that the offset has reached the maximum allowed value.
func ValidateOffset(offset int32) (int32, error) {
	if offset < 0 {
		return 0, fmt.Errorf("offset is zero or less")
	}
	if offset > int32(policy.MaxInt32) {
		return 0, fmt.Errorf("offset reached maximum allowed value")
	}
	return offset, nil
}

// ValidateOrderBy validates the provided sort by field to ensure it meets the required constraints for sorting video retrieval requests. It returns the validated sort by field if it is valid, or an error if it fails any of the validation checks. If the provided sort by field is not valid, the function returns the default value defined by ports.OrderByCreatedAt.
//
// Parameters:
// - orderBy: The string value representing the sort by field to be validated.
//
// Returns:
// - string: The validated sort by field if it is valid, or the default value if the provided value is not valid.
// - error: An error if the sort by field is invalid, with details on the reason for failure.
//
// Validation Rules:
// 1. The sort by field must not exceed the maximum allowed size defined by policy.UrlParamMaxSize. If it does, an error is returned indicating that the sort by field is too large.
// 2. The sort by field must be one of the allowed values defined by ports.OrderByCreatedAt. If it is not, an error is returned indicating that the sort by field is invalid.
func ValidateOrderBy(orderBy string) (string, error) {
	if len(orderBy) > policy.UrlParamMaxSize {
		return "", fmt.Errorf("orderBy parameter is too large")
	}
	switch orderBy {
	case ports.OrderByCreatedAt:
		return orderBy, nil
	default:
		return "", fmt.Errorf("invalid orderBy argument: %s", orderBy)
	}
}

// ValidateIsAsc validates the provided `asc` parameter to determine if sorting should be in ascending order. It returns a boolean value (`true` for ascending, `false` for descending) if the input is valid, or an error if the input fails validation. If the input is invalid, the function returns `false` and an error indicating the issue.
//
// Parameters:
// - asc: The string value representing the sorting direction. It should be either `"t"` (for ascending) or `"f"` (for descending).
//
// Returns:
// - bool: `true` if the sorting direction is ascending, `false` if descending. Returns `false` if the input is invalid (but an error will also be returned in that case).
// - error: An error if the input is invalid, with details on the reason for failure.
//
// Validation Rules:
// 1. The `asc` parameter must not exceed the maximum allowed size defined by policy.UrlParamMaxSize. If it does, an error is returned indicating that the parameter is too large.
// 2. The `asc` parameter must be either `"t"` (for ascending) or `"f"` (for descending). If it is not, an error is returned indicating that the argument is invalid and only `"t"` or `"f"` are allowed.
func ValidateIsAsc(asc string) (bool, error) {
	if len(asc) > policy.UrlParamMaxSize {
		return false, fmt.Errorf("asc parameter is too large")
	}
	switch asc {
	case "t":
		return true, nil
	case "f":
		return false, nil
	default:
		return false, fmt.Errorf("invalid asc argument (only `t` and `f` are allowed)")
	}
}
