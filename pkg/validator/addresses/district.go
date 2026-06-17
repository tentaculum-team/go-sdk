package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type DistrictConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultDistrictConfig() DistrictConfig {
	return DistrictConfig{MaxChars: 100, MinChars: 0, AllErrors: true}
}

// District validates a neighborhood/district name for the global market. Same
// rule set as City: Unicode letters and marks, digits, spaces, hyphens,
// apostrophes, periods and commas.
func District(district string, cfg ...DistrictConfig) error {
	conf := DefaultDistrictConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	district = strings.TrimSpace(district)
	if district == "" {
		return errors.New("district cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(district)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("district cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("district cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(district) || utils.ContainsControlChars(district) {
		errs = append(errs, errors.New("district cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(district) {
		errs = append(errs, errors.New("district contains invalid UTF-8 characters."))
	}
	for _, r := range district {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(" -'.,", r) {
			continue
		}
		errs = append(errs, fmt.Errorf("district contains an invalid character: '%c'.", r))
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
