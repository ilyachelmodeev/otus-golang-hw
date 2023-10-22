package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(stringValue string) (string, error) {
	var buffer string
	result := ""

	for key, value := range stringValue {
		if key == 0 {
			if unicode.IsDigit(value) {
				return "", ErrInvalidString
			}
		}

		if unicode.IsLetter(value) {
			if len(buffer) != 0 {
				result += buffer
				buffer = string(value)
				continue
			}
			buffer += string(value)
			continue
		}

		if unicode.IsDigit(value) {
			if len(buffer) == 0 {
				return "", ErrInvalidString
			}
			result += strings.Repeat(buffer, int(value-'0'))
			buffer = ""
			continue
		}

		return "", ErrInvalidString
	}

	if len(buffer) != 0 {
		result += buffer
	}
	return result, nil
}
