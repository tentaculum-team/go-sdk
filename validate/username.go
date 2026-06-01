package validate

import (
	"errors"
	"fmt"
)

type UsernameConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultUsernameConfig() UsernameConfig {
	return UsernameConfig{
		MaxChars:  50,
		MinChars:  3,
		AllErrors: true,
	}
}

func Username(username string, cfg ...UsernameConfig) error {
	conf := DefaultUsernameConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(username) == 0 {
		message := "The username address cannot be empty."
		errs = append(errs, errors.New(message))
	}

	if len(username) > conf.MaxChars {
		message := fmt.Sprintf(`username cannot exceed %d characters. (current %d)`, conf.MaxChars, len(username))
		errs = append(errs, errors.New(message))
	}

	if len(username) < conf.MinChars {
		message := fmt.Sprintf(`username addresses cannot be shorter than %d characters. (current %d)`, conf.MinChars, len(username))
		errs = append(errs, errors.New(message))
	}

	for i := 0; i < len(username); i++ {
		c := username[i]

		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' ||
			c == '-' ||
			c == '.' {
			continue
		}

		message := fmt.Sprintf("The username contains an invalid character: '%c'.", c)
		errs = append(errs, errors.New(message))
	}

	if isNumericOnly(username) {
		message := "username cannot contain only numbers."
		errs = append(errs, errors.New(message))
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

func isNumericOnly(username string) bool {
	if len(username) == 0 {
		return false
	}

	for i := 0; i < len(username); i++ {
		if username[i] < '0' || username[i] > '9' {
			return false
		}
	}

	return true
}
