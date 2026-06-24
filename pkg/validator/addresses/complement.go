package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type ComplementConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultComplementConfig() ComplementConfig {
	return ComplementConfig{MaxChars: 255, MinChars: 0, AllErrors: true}
}

// Complement validates an address complement/second line for the global market
// ("Apt 4B", "2nd floor", "Bloco C / Sala 12"). Same character set as Street.
func Complement(complement string, cfg ...ComplementConfig) error {
	conf := DefaultComplementConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	complement = strings.TrimSpace(complement)
	if complement == "" {
		return errors.New("complement cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(complement)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("complement cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("complement cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(complement) || utils.ContainsControlChars(complement) {
		errs = append(errs, errors.New("complement cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(complement) {
		errs = append(errs, errors.New("complement contains invalid UTF-8 characters."))
	}
	for _, r := range complement {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(streetAllowedPunct, r) {
			continue
		}
		errs = append(errs, fmt.Errorf("complement contains an invalid character: '%c'.", r))
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
