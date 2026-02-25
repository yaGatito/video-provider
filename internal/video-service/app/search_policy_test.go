package app

import (
	"fmt"
	"testing"
	"video-provider/internal/video-service/policy"
	"video-provider/internal/video-service/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateSearchQuery(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		res, err := ValidateSearchQuery(c.query)
		if c.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Exactly(t, c.outputExpected, res)
		}
	}
}

func TestIncorrectSearchQuery(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{}

	var symbols = "@#$%^&*()+=!?,.;'"
	var format = "se%carch global"

	for i, c := range symbols {
		newCase := struct {
			name  string
			query string
		}{
			name:  fmt.Sprintf("incorrect search: %d; symbol: %c", i+1, c),
			query: fmt.Sprintf(format, c),
		}

		cases = append(cases, newCase)
	}

	for _, c := range cases {
		_, err := ValidateSearchQuery(c.query)
		require.Error(t, err)
	}
}

func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected int32
	}{
		{"negative limit", -5, policy.DefaultVideosLimitPerRequest},
		{"zero limit", 0, policy.DefaultVideosLimitPerRequest},
		{"default limit", policy.DefaultVideosLimitPerRequest, policy.DefaultVideosLimitPerRequest},
		{"above max limit", policy.MaxVideosLimitPerRequest + 1, policy.MaxVideosLimitPerRequest},
		{"valid limit", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateLimit(tt.input)
			require.Exactly(t, tt.expected, result)
		})
	}
}

func TestValidateOffset(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected int32
	}{
		{"negative offset", -5, 0},
		{"zero offset", 0, 0},
		{"positive offset", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOffset(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"invalid sort", "unknown", ports.CreatedAtSort},
		{"valid sort", ports.CreatedAtSort, ports.CreatedAtSort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOrderBy(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateAsc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"ValidEmptyString", "", true}, // default
		{"ValidTrue", "t", true},
		{"ValidFalse", "f", false},
		{"InvalidTrue", "true", false},
		{"ValidFalseString", "false", false},
		{"InvalidNumberOne", "1", false},
		{"InvalidNumberZero", "0", false},
		{"InvalidAsc", "asc", false},
		{"InvalidDesc", "desc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAsc(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateAsc(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
