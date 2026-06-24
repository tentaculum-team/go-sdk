package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type FirstNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultFirstNameConfig() FirstNameConfig {
	return FirstNameConfig{MaxChars: 50, MinChars: 2, AllErrors: true}
}

func FirstName(name string, cfg ...FirstNameConfig) error {
	conf := DefaultFirstNameConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(name) == 0 {
		errs = append(errs, errors.New("first name cannot be empty."))
		return errs[0]
	}

	if utf8.RuneCountInString(name) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("first name cannot exceed %d characters. (current %d)", conf.MaxChars, utf8.RuneCountInString(name)))
	}

	if utf8.RuneCountInString(name) < conf.MinChars {
		errs = append(errs, fmt.Errorf("first name cannot be shorter than %d characters. (current %d)", conf.MinChars, utf8.RuneCountInString(name)))
	}

	if strings.Contains(name, " ") {
		errs = append(errs, errors.New("first name cannot contain spaces."))
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && r != '-' && r != '\'' {
			errs = append(errs, fmt.Errorf("first name contains an invalid character: '%c'.", r))
			break
		}
	}

	if utils.StartsWithInvalidNameChar(name) {
		errs = append(errs, errors.New("first name cannot start with '-' or \"'\"."))
	}

	if utils.EndsWithInvalidNameChar(name) {
		errs = append(errs, errors.New("first name cannot end with '-' or \"'\"."))
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
