package validator

import (
	"errors"
	"fmt"
	"strings"
)

type CountryConfig struct {
	Uppercase bool
	AllErrors bool
}

func DefaultCountryConfig() CountryConfig {
	return CountryConfig{
		Uppercase: true,
		AllErrors: true,
	}
}

// Country validates an ISO 3166-1 alpha-2 country code (BR, US, PT, ...).
func Country(code string, cfg ...CountryConfig) error {
	conf := DefaultCountryConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	code = strings.TrimSpace(code)

	if len(code) == 0 {
		errs = append(errs, errors.New("country code cannot be empty."))
		return errs[0]
	}

	if len(code) != 2 {
		errs = append(errs, fmt.Errorf("country code must be exactly 2 letters (ISO 3166-1 alpha-2). (current %d)", len(code)))
	}

	for i := 0; i < len(code); i++ {
		c := code[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			continue
		}
		errs = append(errs, fmt.Errorf("country code contains an invalid character: '%c'.", c))
		break
	}

	if conf.Uppercase && code != strings.ToUpper(code) {
		errs = append(errs, errors.New("country code must be uppercase."))
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
