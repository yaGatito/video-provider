package policy

import "regexp"

const (
	IDMaxBytes        = 36
	TopicMaxLen       = 48
	TopicMinLen       = 5
	DescriptionMaxLen = 512

	UrlMaxLen    = 100
	SearchMinLen = 3

	ThresholdVideosLimit = 5
	VideosMaxLimit       = 50
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
