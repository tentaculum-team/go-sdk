package validator

import (
	"errors"
	"fmt"
	"strings"
)

type PhoneConfig struct {
	MaxDigits          int
	MinDigits          int
	RequireCountryCode bool
	AllowFormatting    bool
	ReturnAllErrors    bool
}

func DefaultPhoneConfig() PhoneConfig {
	return PhoneConfig{
		MaxDigits:          20,
		MinDigits:          8,
		RequireCountryCode: false,
		AllowFormatting:    true,
		ReturnAllErrors:    true,
	}
}

func Phone(phone string, config ...PhoneConfig) error {
	conf := DefaultPhoneConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	var errs []error
	phone = strings.TrimSpace(phone)

	if len(phone) == 0 {
		errs = append(errs, errors.New("phone cannot be empty."))
		return errs[0]
	}

	if conf.RequireCountryCode && !strings.HasPrefix(phone, "+") {
		errs = append(errs, errors.New("phone must start with the country code '+'."))
	}

	digits := countDigits(phone)

	if digits > conf.MaxDigits {
		errs = append(errs, fmt.Errorf("phone cannot exceed %d digits. (current %d)", conf.MaxDigits, digits))
	}

	if digits < conf.MinDigits {
		errs = append(errs, fmt.Errorf("phone cannot be shorter than %d digits. (current %d)", conf.MinDigits, digits))
	}

	for i := 0; i < len(phone); i++ {

		phoneDigit := phone[i]

		if phoneDigit >= '0' && phoneDigit <= '9' {
			continue
		}
		if phoneDigit == '+' && i == 0 {
			continue
		}
		if !conf.AllowFormatting && (phoneDigit == ' ' || phoneDigit == '-' || phoneDigit == '(' || phoneDigit == ')' || phoneDigit == '.') {
			continue
		}

		errs = append(errs, fmt.Errorf("phone contains an invalid character: '%c'.", phoneDigit))
		break
	}

	if conf.ReturnAllErrors {
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

func countDigits(phoneNumber string) int {
	phoneCount := 0
	for i := 0; i < len(phoneNumber); i++ {
		if phoneNumber[i] >= '0' && phoneNumber[i] <= '9' {
			phoneCount++
		}
	}
	return phoneCount
}
