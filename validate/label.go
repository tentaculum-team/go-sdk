package validate

import (
	"errors"
	"fmt"
)

type LabelConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func Label(label string, cfg LabelConfig) error {
	var errs []error

	if label == "" {
		errs = append(errs, fmt.Errorf("label is required"))
	}

	if cfg.MinChars > 0 && len(label) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("label must be at least %d characters long", cfg.MinChars))
	}

	if cfg.MaxChars > 0 && len(label) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("label must be at most %d characters long", cfg.MaxChars))
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
