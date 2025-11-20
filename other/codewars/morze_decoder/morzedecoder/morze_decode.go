package morzedecoder

// https://www.codewars.com/kata/54b72c16cd7f5154e9000457/train/go
// ···· · −·−−   ·−−− ··− −·· ·
// ···· · -·--    --- ··- -·· ·
// 1100 1100 1100 1100 0000 1100 0000 1111 1100 1100 1111 1100 1111 1100 0000 0000 0000 0011 0011 1111 0011 1111 1001 1111 1100 0000 1100 1100 1111 1100 0000 0111 1110 0110 0110 0000 0011

import (
	"fmt"
	"strings"
)

var MORZE_CODE = map[string]string{
	"-----": "0",
	".----": "1",
	"..---": "2",
	"...--": "3",
	"....-": "4",
	".....": "5",
	"-....": "6",
	"--...": "7",
	"---..": "8",
	"----.": "9",
	".-":    "A", // 10
	"-...":  "B", // 11
	"-.-.":  "C", // 12
	"-..":   "D", // 13
	".":     "E", // 14
	"..-.":  "F", // 15
	"--.":   "G", // 16
	"....":  "H", // 17
	"..":    "I", // 18
	".---":  "J", // 19
	"-.-":   "K", // 20
	".-..":  "L", // 21
	"--":    "M", // 22
	"-.":    "N", // 23
	"---":   "O", // 24
	".--.":  "P", // 25
	"--.-":  "Q", // 26
	".-.":   "R", // 27
	"...":   "S", // 28
	"-":     "T", // 29
	"..-":   "U", // 30
	"...-":  "V", // 31
	".--":   "W", // 32
	"-..-":  "X", // 33
	"-.--":  "Y", // 34
	"--..":  "Z", // 35
}

const (
	ONE_UNIT    = "1"
	ZERO_UNIT   = "0"
	SEVEN_UNITS = "0000000"
	DOT         = "."
	DASH        = "-"
	SPACE       = " "
)

func DecodeBits(bits string) string {
	fmt.Println(bits)
	bits = strings.Trim(bits, ZERO_UNIT)

	determineUnit := func(onesCounter int, zeroesCounter int, threeUnitDot bool) string {
		if (!threeUnitDot && onesCounter >= 1 && onesCounter <= 2) || (threeUnitDot && onesCounter == 3) {
			return DOT
		}
		if onesCounter >= 3 && onesCounter <= 6 {
			return DASH
		}
		if (!threeUnitDot && zeroesCounter >= 3 && zeroesCounter <= 14) || (threeUnitDot && zeroesCounter > 3) {
			return SPACE
		}

		return ""
	}
	sb := strings.Builder{}
	words := strings.Split(bits, SEVEN_UNITS)
	threeUnitDot := false
	if !strings.Contains(bits, "0") && len(bits) > 0 {
		sb.WriteString(DOT)
		return sb.String()
	}
	if len(words) == 1 {
		threeUnitDot = !strings.Contains("0"+words[0]+"0", "010") && !strings.Contains("0"+words[0]+"0", "0110") && len(words[0]) <= 9
	}
	for _, word := range words {
		// word := strings.Split(word, "000")
		onesCounter := 0
		zeroesCounter := 0
		for ixBit, bit := range word {
			if string(bit) == ZERO_UNIT {
				sb.WriteString(determineUnit(onesCounter, 0, threeUnitDot))
				zeroesCounter++
				onesCounter = 0
			} else if string(bit) == ONE_UNIT {
				sb.WriteString(determineUnit(0, zeroesCounter, false))
				onesCounter++
				if ixBit == len(word)-1 {
					sb.WriteString(determineUnit(onesCounter, zeroesCounter, threeUnitDot))
				}
				zeroesCounter = 0
			}
		}

		sb.WriteString(SPACE + SPACE)
	}

	fmt.Println(sb.String())
	return strings.TrimSpace(sb.String())
}

func DecodeMorse(morseCode string) string {
	words := strings.Split(morseCode, "  ")
	sb := strings.Builder{}
	for ixWord, word := range words {
		chars := strings.Split(word, " ")
		for _, char := range chars {
			if char == "" {
				continue
			}
			symbol := MORZE_CODE[char]
			sb.WriteString(symbol)
		}
		if ixWord < len(words)-1 {
			sb.WriteString(" ")
		}
	}

	fmt.Println(sb.String())
	return strings.ReplaceAll(strings.TrimSpace(sb.String()), "  ", " ")
}
