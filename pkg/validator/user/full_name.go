package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type FullNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultFullNameConfig() FullNameConfig {
	return FullNameConfig{MaxChars: 150, MinChars: 5, AllErrors: true}
}

func FullName(name string, cfg ...FullNameConfig) error {
	conf := DefaultFullNameConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(name) == 0 {
		errs = append(errs, errors.New("full name cannot be empty."))
		return errs[0]
	}

	if utf8.RuneCountInString(name) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("full name cannot exceed %d characters. (current %d)", conf.MaxChars, utf8.RuneCountInString(name)))
	}

	if utf8.RuneCountInString(name) < conf.MinChars {
		errs = append(errs, fmt.Errorf("full name cannot be shorter than %d characters. (current %d)", conf.MinChars, utf8.RuneCountInString(name)))
	}

	parts := strings.Fields(name)
	if len(parts) < 2 {
		errs = append(errs, errors.New("full name must contain at least first and last name."))
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && r != '-' && r != '\'' && r != ' ' {
			errs = append(errs, fmt.Errorf("full name contains an invalid character: '%c'.", r))
			break
		}
	}

	if strings.Contains(name, "  ") {
		errs = append(errs, errors.New("full name cannot contain consecutive spaces."))
	}

	if len(name) > 0 && (name[0] == ' ' || utils.StartsWithInvalidNameChar(name)) {
		errs = append(errs, errors.New("full name cannot start with '-', \"'\" or space."))
	}

	if len(name) > 0 && (name[len(name)-1] == ' ' || utils.EndsWithInvalidNameChar(name)) {
		errs = append(errs, errors.New("full name cannot end with '-', \"'\" or space."))
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
