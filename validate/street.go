package validate

import (
	"errors"
	"fmt"
)

type StreetConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func Street(street string, cfg StreetConfig) error {

	var errs []error

	if street == "" {
		errs = append(errs, fmt.Errorf("street is required"))
	}

	if cfg.MinChars > 0 && len(street) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("street must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(street) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("street must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
