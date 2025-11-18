package main

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
	bits = strings.Trim(bits, ZERO_UNIT)
	determineUnit := func(onesCounter int, zeroesCounter int) string {
		if onesCounter >= 1 && onesCounter <= 3 {
			return DOT
		}
		if onesCounter == 6 {
			return DASH
		}
		if zeroesCounter == 6 {
			return SPACE
		}
		return ""
	}
	sb := strings.Builder{}
	words := strings.Split(bits, SEVEN_UNITS)
	for _, word := range words {
		onesCounter := 0
		zeroesCounter := 0
		for ixBit, bit := range word {
			if string(bit) == ZERO_UNIT {
				sb.WriteString(determineUnit(onesCounter, 0))
				zeroesCounter++
				onesCounter = 0
			} else if string(bit) == ONE_UNIT {
				sb.WriteString(determineUnit(0, zeroesCounter))
				onesCounter++
				if ixBit == len(word)-1 {
					sb.WriteString(determineUnit(onesCounter, zeroesCounter))
				}
				zeroesCounter = 0
			}
		}

		sb.WriteString(SPACE)
	}

	return sb.String()
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
	return sb.String()
}

func main() {
	bits := "1100110011001100000011000000111111001100111111001111110000000000000011001111110011111100111111000000110011001111110000001111110011001100000011"
	fmt.Println(DecodeMorse(DecodeBits(bits)))
}
