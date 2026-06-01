package validate

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FirstNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

type LastNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

type FullNameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultFirstNameConfig() FirstNameConfig {
	return FirstNameConfig{MaxChars: 50, MinChars: 2, AllErrors: true}
}

func DefaultLastNameConfig() LastNameConfig {
	return LastNameConfig{MaxChars: 100, MinChars: 2, AllErrors: true}
}

func DefaultFullNameConfig() FullNameConfig {
	return FullNameConfig{MaxChars: 150, MinChars: 5, AllErrors: true}
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

	if startsWithInvalidNameChar(name) {
		errs = append(errs, errors.New("first name cannot start with '-' or \"'\"."))
	}

	if endsWithInvalidNameChar(name) {
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

	if startsWithInvalidNameChar(name) || name[0] == ' ' {
		errs = append(errs, errors.New("last name cannot start with '-', \"'\" or space."))
	}

	if endsWithInvalidNameChar(name) || name[len(name)-1] == ' ' {
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

	if len(name) > 0 && (name[0] == ' ' || startsWithInvalidNameChar(name)) {
		errs = append(errs, errors.New("full name cannot start with '-', \"'\" or space."))
	}

	if len(name) > 0 && (name[len(name)-1] == ' ' || endsWithInvalidNameChar(name)) {
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

func startsWithInvalidNameChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == '-' || s[0] == '\''
}

func endsWithInvalidNameChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s)-1] == '-' || s[len(s)-1] == '\''
}
