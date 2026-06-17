package validate

import (
	"errors"
	"fmt"
)

type CityConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func City(city string, cfg CityConfig) error {
	var errs []error

	if city == "" {
		errs = append(errs, fmt.Errorf("city is required"))
	}

	if cfg.MinChars > 0 && len(city) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("city must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(city) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("city must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
