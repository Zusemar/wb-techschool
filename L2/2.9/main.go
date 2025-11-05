package main

import (
	"errors"
	"strings"
	"unicode"
)

// Unpack unpacks a string like "a4bc2d5e" to "aaaabccddddde"
func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var result strings.Builder
	var prev rune
	var escape bool
	var havePrev bool

	for _, r := range s {
		if escape {
			result.WriteRune(r)
			prev = r
			havePrev = true
			escape = false
		} else if r == '\\' {
			escape = true
		} else if unicode.IsDigit(r) {
			if !havePrev {
				return "", errors.New("invalid string")
			}
			count := int(r - '0')
			for i := 1; i < count; i++ {
				result.WriteRune(prev)
			}
		} else {
			result.WriteRune(r)
			prev = r
			havePrev = true
		}
	}

	if escape {
		return "", errors.New("invalid escape at end")
	}

	return result.String(), nil
}
