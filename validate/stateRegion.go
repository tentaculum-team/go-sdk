package validate

import (
	"errors"
	"fmt"
)

type StateRegionConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func StateRegion(stateRegion string, cfg StateRegionConfig) error {
	var errs []error

	if stateRegion == "" {
		errs = append(errs, fmt.Errorf("state_region is required"))
	}

	if cfg.MinChars > 0 && len(stateRegion) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("state_region must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(stateRegion) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("state_region must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
