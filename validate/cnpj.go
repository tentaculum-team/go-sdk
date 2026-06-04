package validate

import (
	"errors"
	"fmt"
	"strings"
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
	digits := stripNonDigits(cnpj)

	if len(digits) == 0 {
		errs = append(errs, errors.New("CNPJ cannot be empty."))
		return errs[0]
	}

	if len(digits) != 14 {
		errs = append(errs, fmt.Errorf("CNPJ must have 14 digits. (current %d)", len(digits)))
	}

	if len(digits) == 14 && allSameDigit(digits) {
		errs = append(errs, errors.New("CNPJ cannot have all identical digits."))
	}

	if conf.CheckDigits && len(digits) == 14 && !allSameDigit(digits) && !cnpjCheckDigitsValid(digits) {
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

func cnpjCheckDigitsValid(d string) bool {
	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(d[i]-'0') * w1[i]
	}
	r := sum % 11
	dv1 := 0
	if r >= 2 {
		dv1 = 11 - r
	}
	if int(d[12]-'0') != dv1 {
		return false
	}

	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(d[i]-'0') * w2[i]
	}
	r = sum % 11
	dv2 := 0
	if r >= 2 {
		dv2 = 11 - r
	}
	return int(d[13]-'0') == dv2
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func allSameDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}
