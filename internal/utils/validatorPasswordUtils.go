package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var keyboardPatterns = []string{
	"qwerty",
	"asdfgh",
	"zxcvbn",
	"123456",
	"654321",
}

var commonPasswords = map[string]struct{}{
	"password": {},
	"123456":   {},
	"12345678": {},
	"qwerty":   {},
	"admin":    {},
	"senha123": {},
}

func HasNumber(password string) bool {
	for _, r := range password {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func HasLetter(password string) bool {
	for i := 0; i < len(password); i++ {
		c := password[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') {
			return true
		}
	}

	return false
}

func HasSpecialCharacter(password string) bool {
	for i := 0; i < len(password); i++ {
		c := password[i]

		if !(c >= 'a' && c <= 'z') &&
			!(c >= 'A' && c <= 'Z') &&
			!(c >= '0' && c <= '9') {
			return true
		}
	}

	return false
}

func ContainsControlChars(password string) bool {
	for i := 0; i < len(password); i++ {
		if password[i] < 32 || password[i] == 127 {
			return true
		}
	}
	return false
}

func ContainsNullByte(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return true
		}
	}
	return false
}

func ContainsInvalidUTF8(s string) bool {
	return !utf8.ValidString(s)
}

func StartsWithWhitespace(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[0] {
	case ' ', '\t', '\n', '\r':
		return true
	}

	return false
}

func EndsWithWhitespace(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[len(s)-1] {
	case ' ', '\t', '\n', '\r':
		return true
	}

	return false
}

func HasSequentialNumbers(s string) bool {
	count := 1

	for i := 1; i < len(s); i++ {
		prev, curr := s[i-1], s[i]
		if curr >= '0' && curr <= '9' && prev >= '0' && prev <= '9' && curr == prev+1 {
			count++
			if count >= 3 {
				return true
			}
		} else {
			count = 1
		}
	}

	return false
}

func HasSequentialLetters(s string) bool {
	count := 1

	for i := 1; i < len(s); i++ {
		prev, curr := s[i-1], s[i]
		isLetter := func(c byte) bool {
			return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		}
		if isLetter(curr) && isLetter(prev) && curr == prev+1 {
			count++
			if count >= 3 {
				return true
			}
		} else {
			count = 1
		}
	}

	return false
}

func HasKeyboardPatterns(s string) bool {
	s = strings.ToLower(s)

	for _, p := range keyboardPatterns {
		if strings.Contains(s, p) {
			return true
		}
	}

	return false
}

func IsCommonPassword(password string) bool {
	_, exists := commonPasswords[strings.ToLower(password)]
	return exists
}
