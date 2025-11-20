package morzedecoder_test

import (
	"testing"

	"morze/morzedecoder"
)

func TestDecodeBitsTable_Subtests(t *testing.T) {
	cases := []struct {
		name string
		bits string
		want string
	}{
		{"M", "1110111", "M"},
		{"E", "111", "E"},
		{"E", "1111111", "E"},
		{"I", "111000111", "I"},
		// {"HELLO", "111000111000111", "HELLO"},
	}

	for _, tc := range cases {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // optional: run subtests in parallel
			got := morzedecoder.DecodeMorse(morzedecoder.DecodeBits(tc.bits))
			morzeWant := searchOverMap(tc.want)
			morzeGot := searchOverMap(got)
			if got != tc.want {
				t.Fatalf("Decode(%q) = %q [%s]; want %q [%s]", tc.bits, got, morzeGot, tc.want, morzeWant)
			}
		})
	}
}

func searchOverMap(character string) string {
	for code, char := range morzedecoder.MORZE_CODE {
		if char == character {
			return code
		}
	}
	return ""
}
