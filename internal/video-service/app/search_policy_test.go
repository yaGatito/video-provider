package app

import (
	"fmt"
	"testing"
	"video-provider/internal/video-service/policy"

	"github.com/stretchr/testify/require"
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
		{"zero limit", 0, 0, 0, policy.DefaultVideosLimitPerRequest},
		{"negative limit", 5, -1, 5, policy.DefaultVideosLimitPerRequest},
		{"negative offset", -1, 5, 0, 5},
	}

	for _, c := range tests {
		offset, limit := ValidatePagination(c.offset, c.limit)
		require.Exactly(t, c.expectedOffset, offset)
		require.Exactly(t, c.expectedLimit, limit)
	}
}

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
