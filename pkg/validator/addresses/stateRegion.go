package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type StateRegionConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultStateRegionConfig() StateRegionConfig {
	return StateRegionConfig{MaxChars: 100, MinChars: 0, AllErrors: true}
}

// StateRegion validates a state/province/region for the global market. It
// accepts both subdivision codes ("CA", "SP") and full names ("New South
// Wales", "São Paulo"): Unicode letters and marks, digits, spaces, hyphens,
// apostrophes, periods and commas.
func StateRegion(stateRegion string, cfg ...StateRegionConfig) error {
	conf := DefaultStateRegionConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	stateRegion = strings.TrimSpace(stateRegion)
	if stateRegion == "" {
		return errors.New("state_region cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(stateRegion)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("state_region cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("state_region cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(stateRegion) || utils.ContainsControlChars(stateRegion) {
		errs = append(errs, errors.New("state_region cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(stateRegion) {
		errs = append(errs, errors.New("state_region contains invalid UTF-8 characters."))
	}
	for _, r := range stateRegion {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(" -'.,", r) {
			continue
		}
		errs = append(errs, fmt.Errorf("state_region contains an invalid character: '%c'.", r))
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
