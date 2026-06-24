package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type PostalConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultPostalConfig() PostalConfig {
	return PostalConfig{MaxChars: 12, MinChars: 2, AllErrors: true}
}

// Postal validates a postal/ZIP code for the global market. Formats vary widely
// (US "12345-6789", UK "SW1A 1AA", BR "12345-678", NL "1234 AB", JP "100-0001"),
// so it accepts ASCII letters and digits plus spaces and hyphens, bounded by
// length. It does not enforce a country-specific shape.
func Postal(postal string, cfg ...PostalConfig) error {
	conf := DefaultPostalConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	postal = strings.TrimSpace(postal)
	if postal == "" {
		return errors.New("postal code cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(postal)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("postal code cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("postal code cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(postal) || utils.ContainsControlChars(postal) {
		errs = append(errs, errors.New("postal code cannot contain control characters."))
	}
	for _, r := range postal {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || unicode.IsDigit(r) ||
			r == ' ' || r == '-' {
			continue
		}
		errs = append(errs, fmt.Errorf("postal code contains an invalid character: '%c'.", r))
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
