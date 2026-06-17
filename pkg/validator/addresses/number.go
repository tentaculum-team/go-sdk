package validator

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

func DefaultHouseNumberConfig() HouseNumberConfig {
	return HouseNumberConfig{
		MaxChars:    50,
		MinChars:    0,
		AllErrors:   true,
		OnlyNumbers: false,
	}
}

func AddressNumber(number string, cfg ...HouseNumberConfig) error {
	conf := DefaultHouseNumberConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if number == "" {
		errs = append(errs, fmt.Errorf("house number is required"))
	}

	if conf.MinChars > 0 && len(number) < conf.MinChars {
		errs = append(errs, fmt.Errorf("house number is too short"))
	}

	if conf.MaxChars > 0 && len(number) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("house number is too long"))
	}

	if conf.OnlyNumbers {
		for _, char := range number {
			if char < '0' || char > '9' {
				errs = append(errs, fmt.Errorf("house number contains invalid characters"))
				break
			}
		}
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
