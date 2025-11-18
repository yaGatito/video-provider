package main

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

// ···· · −·−−   ·−−− ··− −·· ·
// 1100 1100 1100 1100 0000 1100 0000 1111 1100 1100 1111 1100 1111 1100 0000 0000 0000 0011 0011 1111 0011 1111 1001 1111 1100 0000 1100 1100 1111 1100 0000 0111 1110 0110 0110 0000 0011

var (
	ONE_UNIT    = "1"
	THREE_UNITS = "111"
	ZERO_UNIT   = "0"
	SEVEN_UNITS = "0000000"
	THREE_ZEROS = "000"
	DOT         = "."
	DOT_B       = "1100"
	DASH        = "-"
	DASH_B      = "1111"
	SPACE       = " "
	SPACE_B     = "0000"
)

func DecodeBits(bits string) string {
	bits = strings.Trim(bits, "0")
	sb := strings.Builder{}
	delta := 4
	for i := 0; i < len(bits); i = i + delta {
		if i+delta > len(bits) {
			break
		}
		mrz := bits[i : i+delta]
		switch mrz {
		case DOT_B:
			sb.WriteString(DOT)
		case DASH_B:
			sb.WriteString(DASH)
		case SPACE_B:
			sb.WriteString(SPACE)
		default:
			// do nothing
		}
		fmt.Println(mrz)
	}
	return sb.String()
}

func DecodeMorse(morseCode string) string {
	slice := strings.Split(morseCode, " ")
	sb := strings.Builder{}
	for i := 0; i < len(slice); i++ {
		sb.WriteString(MORZE_CODE[slice[i]])
		sb.WriteString(" ")
	}
	return sb.String()
}

func main() {
	bits := "1100110011001100000011000000111111001100111111001111110000000000000011001111110011111100111111000000110011001111110000001111110011001100000011"
	fmt.Println(DecodeMorse(DecodeBits(bits)))
}
