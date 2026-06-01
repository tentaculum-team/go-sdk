package main

import (
	"fmt"

	"github.com/Tentaculum-dev/go-sdk/validate"
)

func main() {
	var errs []error

	// sem config: usa defaults
	err := validate.Mail("abc")
	errs = append(errs, err)

	// config estilo gin: parte do default, sobrescreve so o que quer
	mailCfg := validate.DefaultEmailConfig()
	mailCfg.MaxChars = 100
	mailCfg.AllowDisposable = true
	err = validate.Mail("abc", mailCfg)
	errs = append(errs, err)

	pwCfg := validate.DefaultPasswordConfig()
	pwCfg.MinChars = 10
	pwCfg.NeedNumbers = true
	pwCfg.NeedLetters = true
	err = validate.Password("abc❤️", pwCfg)
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
