package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tentaculum-team/go-sdk/internal/utils"
)

type LastNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultLastNameConfig() LastNameConfig {
	return LastNameConfig{MaxChars: 100, MinChars: 2, AllErrors: true}
}

func LastName(name string, cfg ...LastNameConfig) error {
	conf := DefaultLastNameConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(name) == 0 {
		errs = append(errs, errors.New("last name cannot be empty."))
		return errs[0]
	}

	if utf8.RuneCountInString(name) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("last name cannot exceed %d characters. (current %d)", conf.MaxChars, utf8.RuneCountInString(name)))
	}

	if utf8.RuneCountInString(name) < conf.MinChars {
		errs = append(errs, fmt.Errorf("last name cannot be shorter than %d characters. (current %d)", conf.MinChars, utf8.RuneCountInString(name)))
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && r != '-' && r != '\'' && r != ' ' {
			errs = append(errs, fmt.Errorf("last name contains an invalid character: '%c'.", r))
			break
		}
	}

	if strings.Contains(name, "  ") {
		errs = append(errs, errors.New("last name cannot contain consecutive spaces."))
	}

	if utils.StartsWithInvalidNameChar(name) || name[0] == ' ' {
		errs = append(errs, errors.New("last name cannot start with '-', \"'\" or space."))
	}

	if utils.EndsWithInvalidNameChar(name) || name[len(name)-1] == ' ' {
		errs = append(errs, errors.New("last name cannot end with '-', \"'\" or space."))
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
