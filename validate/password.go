package validate

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type PasswordConfig struct {
	MaxChars         int
	MinChars         int
	NeedNumbers      bool
	NeedLetters      bool
	NeedSpecialChars bool
	AllErrors        bool
}

func DefaultPasswordConfig() PasswordConfig {
	return PasswordConfig{
		MaxChars:         50,
		MinChars:         6,
		NeedNumbers:      false,
		NeedLetters:      false,
		NeedSpecialChars: true,
		AllErrors:        true,
	}
}

func Password(password string, cfg ...PasswordConfig) error {
	conf := DefaultPasswordConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(password) == 0 {
		message := "The password address cannot be empty."
		errs = append(errs, errors.New(message))
	}

	if len(password) > conf.MaxChars {
		message := fmt.Sprintf(`password cannot exceed %d characters. (current %d)`, conf.MaxChars, len(password))
		errs = append(errs, errors.New(message))
	}

	if len(password) < conf.MinChars {
		message := fmt.Sprintf(`password addresses cannot be shorter than %d characters. (current %d)`, conf.MinChars, len(password))
		errs = append(errs, errors.New(message))
	}

	if conf.NeedNumbers {
		if !hasNumber(password) {
			message := "The password must contain at least one number."
			errs = append(errs, errors.New(message))
		}
	}

	if conf.NeedLetters {
		if !hasLetter(password) {
			message := "The password must contain at least one letter."
			errs = append(errs, errors.New(message))
		}
	}

	if conf.NeedSpecialChars {
		if !hasSpecialCharacter(password) {
			message := "The password needs a special character."
			errs = append(errs, errors.New(message))
		}
	}

	if containsControlChars(password) {
		message := "The password cannot contain control characters."
		errs = append(errs, errors.New(message))
	}

	if containsNullByte(password) {
		message := "The password cannot contain null bytes."
		errs = append(errs, errors.New(message))
	}

	if containsInvalidUTF8(password) {
		message := "The password contains invalid UTF-8 characters."
		errs = append(errs, errors.New(message))
	}

	if startsWithWhitespace(password) {
		message := "The password cannot start with whitespace."
		errs = append(errs, errors.New(message))
	}

	if endsWithWhitespace(password) {
		message := "The password cannot end with whitespace."
		errs = append(errs, errors.New(message))
	}

	if hasSequentialNumbers(password) {
		message := "The password cannot contain sequential numbers."
		errs = append(errs, errors.New(message))
	}

	if hasSequentialLetters(password) {
		message := "The password cannot contain sequential letters."
		errs = append(errs, errors.New(message))
	}

	if hasKeyboardPatterns(password) {
		message := "The password cannot contain common keyboard patterns."
		errs = append(errs, errors.New(message))
	}

	if isCommonPassword(password) {
		message := "The password is too common."
		errs = append(errs, errors.New(message))
	}

	if conf.AllErrors {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	} else {
		if len(errs) > 0 {
			return errs[0]
		}
	}

	return nil
}

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

func hasNumber(password string) bool {
	for _, r := range password {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasLetter(password string) bool {
	for i := 0; i < len(password); i++ {
		c := password[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') {
			return true
		}
	}

	return false
}

func hasSpecialCharacter(password string) bool {
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

func containsControlChars(password string) bool {
	for i := 0; i < len(password); i++ {
		if password[i] < 32 || password[i] == 127 {
			return true
		}
	}
	return false
}

func containsNullByte(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return true
		}
	}
	return false
}

func containsInvalidUTF8(s string) bool {
	return !utf8.ValidString(s)
}

func startsWithWhitespace(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[0] {
	case ' ', '\t', '\n', '\r':
		return true
	}

	return false
}

func endsWithWhitespace(s string) bool {
	if len(s) == 0 {
		return false
	}

	switch s[len(s)-1] {
	case ' ', '\t', '\n', '\r':
		return true
	}

	return false
}

func hasSequentialNumbers(s string) bool {
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

func hasSequentialLetters(s string) bool {
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

func hasKeyboardPatterns(s string) bool {
	s = strings.ToLower(s)

	for _, p := range keyboardPatterns {
		if strings.Contains(s, p) {
			return true
		}
	}

	return false
}

func isCommonPassword(password string) bool {
	_, exists := commonPasswords[strings.ToLower(password)]
	return exists
}
