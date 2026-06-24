package validator

import (
	"errors"
	"fmt"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
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
		usernameChar := username[i]

		if (usernameChar >= 'a' && usernameChar <= 'z') ||
			(usernameChar >= 'A' && usernameChar <= 'Z') ||
			(usernameChar >= '0' && usernameChar <= '9') ||
			usernameChar == '_' ||
			usernameChar == '-' ||
			usernameChar == '.' {
			continue
		}

		message := fmt.Sprintf("The username contains an invalid character: '%c'.", usernameChar)
		errs = append(errs, errors.New(message))
	}

	if utils.IsNumericOnly(username) {
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
