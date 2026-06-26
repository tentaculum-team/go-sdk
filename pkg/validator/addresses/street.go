package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tentaculum-team/go-sdk/internal/utils"
)

// streetAllowedPunct are the non-alphanumeric characters permitted in street
// and complement lines worldwide (e.g. "Av. Paulista, 1000", "5ª/Rua-B #2").
const streetAllowedPunct = " .,'-/#°ºª()&"

type StreetConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultStreetConfig() StreetConfig {
	return StreetConfig{MaxChars: 255, MinChars: 0, AllErrors: true}
}

// Street validates a street/address line for the global market. It accepts
// Unicode letters and marks, digits, spaces and common address punctuation,
// rejecting control characters, null bytes and invalid UTF-8.
func Street(street string, cfg ...StreetConfig) error {
	conf := DefaultStreetConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	street = strings.TrimSpace(street)
	if street == "" {
		return errors.New("street cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(street)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("street cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("street cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(street) || utils.ContainsControlChars(street) {
		errs = append(errs, errors.New("street cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(street) {
		errs = append(errs, errors.New("street contains invalid UTF-8 characters."))
	}
	for _, r := range street {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(streetAllowedPunct, r) {
			continue
		}
		errs = append(errs, fmt.Errorf("street contains an invalid character: '%c'.", r))
		break
	}

	if conf.AllErrors {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	} else if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
