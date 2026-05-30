package validate

import (
	"errors"
	"fmt"
	"strings"
)

type EmailConfig struct {
	MaxChars        *int
	MinChars        *int
	AllowDisposable bool
	AllErrors       bool
}

func Mail(mail string, cfg ...EmailConfig) error {
	var conf EmailConfig
	if len(cfg) > 0 {
		conf = cfg[0]
	} else {
		conf.AllErrors = true
	}

	var errs []error
	mail = strings.TrimSpace(mail)

	if len(mail) == 0 {
		message := "The email address cannot be empty."
		errs = append(errs, errors.New(message))
	}

	if conf.MaxChars == nil {
		if len(mail) > 200 {
			message := fmt.Sprintf(`mails cannot exceed 200 characters. (current %d)`, len(mail))
			errs = append(errs, errors.New(message))
		}
	} else {
		if len(mail) > *conf.MaxChars {
			message := fmt.Sprintf(`mails cannot exceed %d characters. (current %d)`, *conf.MaxChars, len(mail))
			errs = append(errs, errors.New(message))
		}
	}

	if conf.MinChars == nil {
		if len(mail) < 6 {
			message := fmt.Sprintf(`mail addresses cannot be shorter than 6 characters. (current %d)`, len(mail))
			errs = append(errs, errors.New(message))
		}
	} else {
		if len(mail) < *conf.MinChars {
			message := fmt.Sprintf(`mail addresses cannot be shorter than %d characters. (current %d)`, *conf.MinChars, len(mail))
			errs = append(errs, errors.New(message))
		}
	}

	if !strings.Contains(mail, "@") {
		message := fmt.Sprintf(`This is not an mail, because it is missing "%s"`, `@`)
		errs = append(errs, errors.New(message))
	}

	if !strings.Contains(mail, ".") {
		message := fmt.Sprintf(`This is not an mail, because it is missing "%s"`, `.`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, " ") {
		message := `mails cannot contain spaces.`
		errs = append(errs, errors.New(message))
	}

	if hasAccent(mail) {
		message := `mail cannot contain accents`
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "<") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `<`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, ">") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `>`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "(") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `(`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, ")") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `)`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "[") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `[`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "]") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `]`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, ",") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `,`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, ";") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `;`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, ":") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `:`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "\\") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `\\`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "/") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `/`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "\"") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `\`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "'") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `'`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "!") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `!`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "#") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `#`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "$") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `$`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "%") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `%`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "^") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `^`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "&") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `&`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "*") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `*`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "=") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `=`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "+") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `+`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "{") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `{`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "}") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `}`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "|") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `|`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "?") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `?`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "~") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, `~`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "`") {
		message := fmt.Sprintf(`mails cannot contain "%s"`, "`")
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "' #") {
		message := "mail format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "'/*") {
		message := "mail format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "' or true --") {
		message := "mail format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(mail, "\") or (\"1\"=\"1") {
		message := "mail format is invalid"
		errs = append(errs, errors.New(message))
	}

	if mail != strings.ToLower(mail) {
		message := "mail addresses cannot be in uppercase."
		errs = append(errs, errors.New(message))
	}

	if !conf.AllowDisposable {
		disposable, err := isDisposable(mail)
		if err != nil {
			errs = append(errs, errors.New("internal error"))
		} else if disposable {
			errs = append(errs, errors.New("temporary mails are not allowed"))
		}
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
