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

// ^[\p{L}\p{N}\s_\-!?;.,]{2,100}$
// ^[\p{L}\p{N}]]{2,}

func GET_TEXTING_FORMAT_RE_128() *regexp.Regexp {
	return regexp.MustCompile(`^[\pL\pN\s]{2}[_\-!?;.,\pL\pN\s]{0,126}$`)
	// return regexp.MustCompile(`^[\p{L}]{2,128}[_\-!?;.,\p{L}\p{N}\s]{,126}$`)
}

func GET_LARGE_TEXT_FORMAT_RE_512() *regexp.Regexp {
	return regexp.MustCompile(`^^[\pL\pN\s]{2}[_\-!?;.,\pL\pN\s]{0,510}$`)
}

func GET_WORDS_FORMAT_RE_128() *regexp.Regexp {
	return regexp.MustCompile(`^[\p{L}\p{N}\s]{2,128}$`)
}

// var GET_TEXTING_FORMAT_RE_128 = regexp.MustCompile(`^[_\-!?;.,\p{L}\p{N}\s]{2,512}$`)
// var LargeTextFormatRe512 = regexp.MustCompile(`^[_\-!?;.,\p{L}\p{N}\s]{2,512}$`)
// var WordsFormatRe128 = regexp.MustCompile(`^[\p{L}\p{N}\s]{1,100}$`)

//  aa !!! asdasd !!!! asdsad s!!!!  ads ada  !!!! dwadsa d adasd  12 we1  2212 123 21 13  12321 213 ! ! 2321  21 321321 213 21!  1 ! 23  31 12 12321! 1 231 23 32  !! ! 1!  1231 32132 ! ! 12 121 2 13231321 3 21321321 313211322!!
