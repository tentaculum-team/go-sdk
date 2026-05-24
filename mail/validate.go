package mail

import (
	"errors"
	"fmt"
)

func Validate(email string) error {

	var errs []error

	if len(email) > 255 {
		message := fmt.Sprintf("Emails cannot exceed 255 characters. (current %d)", len(email))
		errs = append(errs, errors.New(message))
	}

	if len(email) < 6 {
		message := fmt.Sprintf("Email addresses cannot be shorter than 10 characters. (current %d)", len(email))
		errs = append(errs, errors.New(message))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
