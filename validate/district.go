package validate

import (
	"errors"
	"fmt"
)

type DistrictConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func District(district string, cfg DistrictConfig) error {
	var errs []error

	if district == "" {
		errs = append(errs, fmt.Errorf("district is required"))
	}

	if cfg.MinChars > 0 && len(district) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("district must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(district) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("district must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
