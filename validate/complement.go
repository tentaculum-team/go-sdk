package validate

import (
	"errors"
	"fmt"
)

type ComplementConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func Complement(complement string, cfg ComplementConfig) error {
	var errs []error
	if complement == "" {
		errs = append(errs, fmt.Errorf("complement is required"))
	}

	if cfg.MinChars > 0 && len(complement) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("complement must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(complement) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("complement must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
