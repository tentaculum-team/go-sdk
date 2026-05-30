package main

import (
	"fmt"

	"github.com/ViitoJooj/sdk/validate"
)

func main() {
	var errs []error

	err := validate.Mail("abc")
	errs = append(errs, err)
	err = validate.Password("abc❤️")
	errs = append(errs, err)
	err = validate.Username("abc❤️")
	errs = append(errs, err)
	err = validate.FirstName("João❤️")
	errs = append(errs, err)
	err = validate.LastName("Vitor❤️")
	errs = append(errs, err)
	err = validate.FullName("João Vitor ❤️Santana Oqueres")
	errs = append(errs, err)

	for i := 0; i < len(errs); i++ {
		fmt.Println(errs[i])
	}

}
