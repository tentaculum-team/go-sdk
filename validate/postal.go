package validate

import (
	"errors"
	"fmt"
	"strings"
)

type PostalConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultPostalConfig() PostalConfig {
	return PostalConfig{
		MaxChars:  20,
		MinChars:  4,
		AllErrors: true,
	}
}

func Postal(postal string, cfg ...PostalConfig) error {
	conf := DefaultPostalConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	postal = strings.TrimSpace(postal)

	if len(postal) == 0 {
		errs = append(errs, errors.New("postal code cannot be empty."))
		return errs[0]
	}

	if len(postal) > conf.MaxChars {
		errs = append(errs, fmt.Errorf("postal code cannot exceed %d characters. (current %d)", conf.MaxChars, len(postal)))
	}

	if len(postal) < conf.MinChars {
		errs = append(errs, fmt.Errorf("postal code cannot be shorter than %d characters. (current %d)", conf.MinChars, len(postal)))
	}

	for i := 0; i < len(postal); i++ {
		c := postal[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == ' ' ||
			c == '-' {
			continue
		}
		errs = append(errs, fmt.Errorf("postal code contains an invalid character: '%c'.", c))
		break
	}

	if conf.AllErrors {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	} else {
		if len(errs) > 0 {
			return errs[0]
		}
	}

	return nil
}
