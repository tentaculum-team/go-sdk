package validator

import (
	"errors"
	"fmt"
)

type DistrictConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultDistrictConfig() DistrictConfig {
	return DistrictConfig{
		MaxChars:  100,
		MinChars:  0,
		AllErrors: true,
	}
}

func District(district string, cfg ...DistrictConfig) error {
	conf := DefaultDistrictConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if district == "" {
		errs = append(errs, fmt.Errorf("district is required"))
	}

	if conf.MinChars > 0 && len(district) < conf.MinChars {
		errs = append(errs, fmt.Errorf("district must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(district) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("district must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
