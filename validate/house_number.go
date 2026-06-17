package validate

import (
	"errors"
	"fmt"
)

type HouseNumberConfig struct {
	MaxChars    int
	MinChars    int
	AllErrors   bool
	OnlyNumbers bool
}

func HouseNumber(number string, cfg HouseNumberConfig) error {
	var errs []error
	if number == "" {
		errs = append(errs, fmt.Errorf("house number is required"))
	}

	if cfg.MinChars > 0 && len(number) < cfg.MinChars {
		errs = append(errs, fmt.Errorf("house number is too short"))
	}

	if cfg.MaxChars > 0 && len(number) > cfg.MaxChars {
		errs = append(errs, fmt.Errorf("house number is too long"))
	}

	if cfg.OnlyNumbers {
		for _, char := range number {
			if char < '0' || char > '9' {
				errs = append(errs, fmt.Errorf("house number contains invalid characters"))
				break
			}
		}
	}

	if len(errs) > 0 {
		if cfg.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
