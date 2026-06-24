package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type HouseNumberConfig struct {
	MaxChars    int
	MinChars    int
	AllErrors   bool
	OnlyNumbers bool
}

func DefaultHouseNumberConfig() HouseNumberConfig {
	return HouseNumberConfig{MaxChars: 50, MinChars: 0, AllErrors: true, OnlyNumbers: false}
}

// AddressNumber validates a house/building number for the global market. House
// numbers are alphanumeric worldwide ("221B", "12/3", "1-3", "s/n", "Plot 5"),
// so by default letters and the separators "/-. " are allowed. Set OnlyNumbers
// to restrict to digits.
func AddressNumber(number string, cfg ...HouseNumberConfig) error {
	conf := DefaultHouseNumberConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	number = strings.TrimSpace(number)
	if number == "" {
		return errors.New("house number cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(number)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("house number cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("house number cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(number) || utils.ContainsControlChars(number) {
		errs = append(errs, errors.New("house number cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(number) {
		errs = append(errs, errors.New("house number contains invalid UTF-8 characters."))
	}
	for _, r := range number {
		if conf.OnlyNumbers {
			if unicode.IsDigit(r) {
				continue
			}
			errs = append(errs, errors.New("house number must contain only digits."))
			break
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(" /-.°ºª", r) {
			continue
		}
		errs = append(errs, fmt.Errorf("house number contains an invalid character: '%c'.", r))
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
