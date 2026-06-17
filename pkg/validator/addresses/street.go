package validator

import (
	"errors"
	"fmt"
)

type StreetConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultStreetConfig() StreetConfig {
	return StreetConfig{
		MaxChars:  255,
		MinChars:  0,
		AllErrors: true,
	}
}

func Street(street string, cfg ...StreetConfig) error {
	conf := DefaultStreetConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if street == "" {
		errs = append(errs, fmt.Errorf("street is required"))
	}

	if conf.MinChars > 0 && len(street) < conf.MinChars {
		errs = append(errs, fmt.Errorf("street must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(street) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("street must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
