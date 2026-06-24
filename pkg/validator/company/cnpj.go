package validator

import (
	"errors"
	"fmt"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type CNPJConfig struct {
	CheckDigits bool
	AllErrors   bool
}

func DefaultCNPJConfig() CNPJConfig {
	return CNPJConfig{
		CheckDigits: true,
		AllErrors:   true,
	}
}

func CNPJ(cnpj string, cfg ...CNPJConfig) error {
	conf := DefaultCNPJConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	digits := utils.StripNonDigits(cnpj)

	if len(digits) == 0 {
		errs = append(errs, errors.New("CNPJ cannot be empty."))
		return errs[0]
	}

	if len(digits) != 14 {
		errs = append(errs, fmt.Errorf("CNPJ must have 14 digits. (current %d)", len(digits)))
	}

	if len(digits) == 14 && utils.AllSameDigit(digits) {
		errs = append(errs, errors.New("CNPJ cannot have all identical digits."))
	}

	if conf.CheckDigits && len(digits) == 14 && !utils.AllSameDigit(digits) && !utils.CNPJCheckDigitsValid(digits) {
		errs = append(errs, errors.New("CNPJ check digits are invalid."))
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
