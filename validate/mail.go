package validate

import "errors"

func Email(email string) error {
	if len(email) > 50 {
		return errors.New("email cannot be must 50")
	}

	return nil
}
