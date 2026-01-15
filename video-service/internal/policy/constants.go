package policy

import "regexp"

const (
	MAX_ID_BYTES_SIZE          = 36
	MAX_TOPIC_BYTES_SIZE       = 48
	MAX_DESCRIPTION_BYTES_SIZE = 512

	MAX_SEARCH_BYTES_SIZE        = 100
	MIN_SEARCH_BYTES_SIZE        = 2
	MAX_VIDEOS_LIMIT_PER_REQUEST = 50
)

// var TextingFormatRe = regexp.MustCompile(`^[\p{L}\p{N}\s_-!?;.,]{1,100}$`)
var WordsFormatRe = regexp.MustCompile(`^[\p{L}\p{N}\s]{1,100}$`)
