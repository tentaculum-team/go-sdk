package validate

import (
	"errors"
	"fmt"
	"strings"
)

type PhoneConfig struct {
	MaxChars    int // max number of digits
	MinChars    int // min number of digits
	RequirePlus bool
	DigitsOnly  bool
	AllErrors   bool
}

func DefaultPhoneConfig() PhoneConfig {
	return PhoneConfig{
		MaxChars:    20,
		MinChars:    8,
		RequirePlus: false,
		DigitsOnly:  true,
		AllErrors:   true,
	}
}

func Phone(phone string, cfg ...PhoneConfig) error {
	conf := DefaultPhoneConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	phone = strings.TrimSpace(phone)

	if len(phone) == 0 {
		errs = append(errs, errors.New("phone cannot be empty."))
		return errs[0]
	}

	if conf.RequirePlus && !strings.HasPrefix(phone, "+") {
		errs = append(errs, errors.New("phone must start with the country code '+'."))
	}

	digits := countDigits(phone)

	if digits > conf.MaxChars {
		errs = append(errs, fmt.Errorf("phone cannot exceed %d digits. (current %d)", conf.MaxChars, digits))
	}

	if digits < conf.MinChars {
		errs = append(errs, fmt.Errorf("phone cannot be shorter than %d digits. (current %d)", conf.MinChars, digits))
	}

	for i := 0; i < len(phone); i++ {
		c := phone[i]
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '+' && i == 0 {
			continue
		}
		if !conf.DigitsOnly && (c == ' ' || c == '-' || c == '(' || c == ')' || c == '.') {
			continue
		}
		errs = append(errs, fmt.Errorf("phone contains an invalid character: '%c'.", c))
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

func countDigits(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n++
		}
	}
	return n
}
