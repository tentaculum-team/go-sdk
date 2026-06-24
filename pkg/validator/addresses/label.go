package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type LabelConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultLabelConfig() LabelConfig {
	return LabelConfig{MaxChars: 50, MinChars: 0, AllErrors: true}
}

// Label validates a user-chosen address label ("Home", "Casa", "Work #2") for
// the global market: Unicode letters and marks, digits, spaces and a small set
// of punctuation.
func Label(label string, cfg ...LabelConfig) error {
	conf := DefaultLabelConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	label = strings.TrimSpace(label)
	if label == "" {
		return errors.New("label cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(label)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("label cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("label cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(label) || utils.ContainsControlChars(label) {
		errs = append(errs, errors.New("label cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(label) {
		errs = append(errs, errors.New("label contains invalid UTF-8 characters."))
	}
	for _, r := range label {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(" .,'-_#", r) {
			continue
		}
		errs = append(errs, fmt.Errorf("label contains an invalid character: '%c'.", r))
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
