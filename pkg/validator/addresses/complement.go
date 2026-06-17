package validator

import (
	"errors"
	"fmt"
)

type ComplementConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultComplementConfig() ComplementConfig {
	return ComplementConfig{
		MaxChars:  255,
		MinChars:  0,
		AllErrors: true,
	}
}

func Complement(complement string, cfg ...ComplementConfig) error {
	conf := DefaultComplementConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if complement == "" {
		errs = append(errs, fmt.Errorf("complement is required"))
	}

	if conf.MinChars > 0 && len(complement) < conf.MinChars {
		errs = append(errs, fmt.Errorf("complement must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(complement) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("complement must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
