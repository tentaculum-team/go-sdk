package validator

import (
	"errors"
	"fmt"

	"github.com/tentaculum-team/go-sdk/internal/utils"
)

type CPFConfig struct {
	CheckDigits bool
	AllErrors   bool
}

func DefaultCPFConfig() CPFConfig {
	return CPFConfig{
		CheckDigits: true,
		AllErrors:   true,
	}
}

func CPF(cpf string, cfg ...CPFConfig) error {
	conf := DefaultCPFConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	digits := utils.StripNonDigits(cpf)

	if len(digits) == 0 {
		errs = append(errs, errors.New("CPF cannot be empty."))
		return errs[0]
	}

	if len(digits) != 11 {
		errs = append(errs, fmt.Errorf("CPF must have 11 digits. (current %d)", len(digits)))
	}

	if len(digits) == 11 && utils.AllSameDigit(digits) {
		errs = append(errs, errors.New("CPF cannot have all identical digits."))
	}

	if conf.CheckDigits &&
		len(digits) == 11 &&
		!utils.AllSameDigit(digits) &&
		!utils.CPFCheckDigitsValid(digits) {
		errs = append(errs, errors.New("CPF check digits are invalid."))
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
