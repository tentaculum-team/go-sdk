package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type EmailConfig struct {
	MaxChars        int
	MinChars        int
	AllowDisposable bool
	AllErrors       bool
}

func DefaultEmailConfig() EmailConfig {
	return EmailConfig{
		MaxChars:        150,
		MinChars:        6,
		AllowDisposable: false,
		AllErrors:       true,
	}
}

func hasAccent(mail string) bool {
	for _, r := range mail {
		if r > 127 {
			return true
		}
	}
	return false
}

func Mail(mail string, cfg ...EmailConfig) error {
	conf := DefaultEmailConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	var errs []error
	mail = strings.TrimSpace(mail)

	if len(mail) == 0 {
		message := "The email address cannot be empty."
		errs = append(errs, errors.New(message))
	}

	if len(mail) > conf.MaxChars {
		message := fmt.Sprintf(`mails cannot exceed %d characters. (current %d)`, conf.MaxChars, len(mail))
		errs = append(errs, errors.New(message))
	}

	if len(mail) < conf.MinChars {
		message := fmt.Sprintf(`mail addresses cannot be shorter than %d characters. (current %d)`, conf.MinChars, len(mail))
		errs = append(errs, errors.New(message))
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
		if disposable, err := utils.IsDisposable(mail); err == nil && disposable {
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
