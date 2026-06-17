package validator

import (
	"errors"
	"fmt"
)

type StateRegionConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultStateRegionConfig() StateRegionConfig {
	return StateRegionConfig{
		MaxChars:  100,
		MinChars:  0,
		AllErrors: true,
	}
}

func StateRegion(stateRegion string, cfg ...StateRegionConfig) error {
	conf := DefaultStateRegionConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if stateRegion == "" {
		errs = append(errs, fmt.Errorf("state_region is required"))
	}

	if conf.MinChars > 0 && len(stateRegion) < conf.MinChars {
		errs = append(errs, fmt.Errorf("state_region must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(stateRegion) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("state_region must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
