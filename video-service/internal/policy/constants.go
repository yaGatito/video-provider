package policy

import "regexp"

const (
	MaxIDBytesSize          = 36
	MaxTopicBytesSize       = 48
	MaxDescriptionBytesSize = 512

	MaxSearchBytesSize           = 100
	MinSearchBytesSize           = 2
	DefaultVideosLimitPerRequest = 5
	MaxVideosLimitPerRequest     = 50
)

func GetTextingFormateRE128() *regexp.Regexp {
	return regexp.MustCompile(`^[\pL\pN\s]{2}[_\-!?;.,\pL\pN\s]{0,126}$`)
}

func GetLargeTextFormatRE512() *regexp.Regexp {
	return regexp.MustCompile(`^^[\pL\pN\s]{2}[_\-!?;.,\pL\pN\s]{0,510}$`)
}

func GetWordsFormatRE128() *regexp.Regexp {
	return regexp.MustCompile(`^[\p{L}\p{N}\s]{2,128}$`)
}
