package httpadp

import (
	"fmt"
	"testing"
	"video-provider/internal/video-service/policy"

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
		res, err := validateSearchQuery(c.query)
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
		_, err := validateSearchQuery(c.query)
		require.Error(t, err)
	}
}

func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		input   int32
	}{
		{"ok", false, 5},
		{"negative limit", true, -5},
		{"zero limit", true, 0},
		{"above max limit", true, policy.VideosMaxLimit + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateLimit(tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Exactly(t, tt.input, result)
			}
		})
	}
}

func TestValidateOffset(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		input   int32
	}{
		{"positive offset", false, 10},
		{"zero offset", false, 0},
		{"negative offset", true, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateOffset(tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Exactly(t, tt.input, result)
			}
		})
	}
}

func TestValidateOrderBy(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		input   string
	}{
		{"valid sort", false, "createdAt"},
		{"invalid sort", true, "CreateddAte"},
	}

	for _, tt := range tests {
		result, err := validateOrderBy(tt.input)

		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Exactly(t, tt.input, result)
		}
	}
}

func TestValidateAsc(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     bool
		input       string
		expectedRes bool
	}{
		{"ok true (asc)", false, "t", true},
		{"ok false (desc)", false, "f", false},
		{"invalid value", true, "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateIsAsc(tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Exactly(t, tt.expectedRes, result)
			}
		})
	}
}
