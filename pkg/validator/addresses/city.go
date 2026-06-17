package validator

import (
	"errors"
	"fmt"
)

type CityConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultCityConfig() CityConfig {
	return CityConfig{
		MaxChars:  100,
		MinChars:  0,
		AllErrors: true,
	}
}

func City(city string, cfg ...CityConfig) error {
	conf := DefaultCityConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if city == "" {
		errs = append(errs, fmt.Errorf("city is required"))
	}

	if conf.MinChars > 0 && len(city) < conf.MinChars {
		errs = append(errs, fmt.Errorf("city must be at least %d characters long", conf.MinChars))
	}

	if conf.MaxChars > 0 && len(city) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("city must be at most %d characters long", conf.MaxChars))
	}

	if len(errs) > 0 {
		if conf.AllErrors {
			return errors.Join(errs...)
		}
		return errs[0]
	}

	return nil
}
