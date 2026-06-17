package validator

import (
	"errors"
	"fmt"
)

type LabelConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		MaxChars:  50,
		MinChars:  0,
		AllErrors: true,
	}
}

func Label(label string, cfg ...LabelConfig) error {
	conf := DefaultLabelConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if label == "" {
		errs = append(errs, fmt.Errorf("label is required"))
	}

	if conf.MinChars > 0 && len(label) < conf.MinChars {
		errs = append(errs, fmt.Errorf("label must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(label) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("label must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
