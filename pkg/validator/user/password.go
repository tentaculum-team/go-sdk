package validator

import (
	"errors"
	"fmt"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type PasswordConfig struct {
	MaxChars         int
	MinChars         int
	NeedNumbers      bool
	NeedLetters      bool
	NeedSpecialChars bool
	AllErrors        bool
}

func DefaultPasswordConfig() PasswordConfig {
	return PasswordConfig{
		MaxChars:         50,
		MinChars:         6,
		NeedNumbers:      false,
		NeedLetters:      false,
		NeedSpecialChars: true,
		AllErrors:        true,
	}
}

func Password(password string, cfg ...PasswordConfig) error {
	conf := DefaultPasswordConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error

	if len(password) == 0 {
		message := "The password address cannot be empty."
		errs = append(errs, errors.New(message))
	}

	if len(password) > conf.MaxChars {
		message := fmt.Sprintf(`password cannot exceed %d characters. (current %d)`, conf.MaxChars, len(password))
		errs = append(errs, errors.New(message))
	}

	if len(password) < conf.MinChars {
		message := fmt.Sprintf(`password addresses cannot be shorter than %d characters. (current %d)`, conf.MinChars, len(password))
		errs = append(errs, errors.New(message))
	}

	if conf.NeedNumbers {
		if !utils.HasNumber(password) {
			message := "The password must contain at least one number."
			errs = append(errs, errors.New(message))
		}
	}

	if conf.NeedLetters {
		if !utils.HasLetter(password) {
			message := "The password must contain at least one letter."
			errs = append(errs, errors.New(message))
		}
	}

	if conf.NeedSpecialChars {
		if !utils.HasSpecialCharacter(password) {
			message := "The password needs a special character."
			errs = append(errs, errors.New(message))
		}
	}

	if utils.ContainsControlChars(password) {
		message := "The password cannot contain control characters."
		errs = append(errs, errors.New(message))
	}

	if utils.ContainsNullByte(password) {
		message := "The password cannot contain null bytes."
		errs = append(errs, errors.New(message))
	}

	if utils.ContainsInvalidUTF8(password) {
		message := "The password contains invalid UTF-8 characters."
		errs = append(errs, errors.New(message))
	}

	if utils.StartsWithWhitespace(password) {
		message := "The password cannot start with whitespace."
		errs = append(errs, errors.New(message))
	}

	if utils.EndsWithWhitespace(password) {
		message := "The password cannot end with whitespace."
		errs = append(errs, errors.New(message))
	}

	if utils.HasSequentialNumbers(password) {
		message := "The password cannot contain sequential numbers."
		errs = append(errs, errors.New(message))
	}

	if utils.HasSequentialLetters(password) {
		message := "The password cannot contain sequential letters."
		errs = append(errs, errors.New(message))
	}

	if utils.HasKeyboardPatterns(password) {
		message := "The password cannot contain common keyboard patterns."
		errs = append(errs, errors.New(message))
	}

	if utils.IsCommonPassword(password) {
		message := "The password is too common."
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
