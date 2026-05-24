package mail

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	domains map[string]struct{}
	once    sync.Once
	loadErr error
)

func Validate(email string) error {

	var errs []error
	email = strings.TrimSpace(email)

	if len(email) > 255 {
		message := fmt.Sprintf(`Emails cannot exceed 255 characters. (current %d)`, len(email))
		errs = append(errs, errors.New(message))
	}

	if len(email) < 6 {
		message := fmt.Sprintf(`Email addresses cannot be shorter than 10 characters. (current %d)`, len(email))
		errs = append(errs, errors.New(message))
	}

	if !strings.Contains(email, "@") {
		message := fmt.Sprintf(`This is not an email, because it is missing "%s"`, `@`)
		errs = append(errs, errors.New(message))
	}

	if !strings.Contains(email, ".") {
		message := fmt.Sprintf(`This is not an email, because it is missing "%s"`, `.`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, " ") {
		message := `Emails cannot contain spaces.`
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "<") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `<`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, ">") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `>`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "(") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `(`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, ")") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `)`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "[") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `[`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "]") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `]`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, ",") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `,`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, ";") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `;`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, ":") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `:`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "\\") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `\\`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "/") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `/`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "\"") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `\`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "'") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `'`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "!") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `!`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "#") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `#`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "$") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `$`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "%") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `%`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "^") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `^`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "&") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `&`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "*") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `*`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "=") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `=`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "+") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `+`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "{") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `{`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "}") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `}`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "|") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `|`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "?") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `?`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "~") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, `~`)
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "`") {
		message := fmt.Sprintf(`Emails cannot contain "%s"`, "`")
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "' or 1=1 --") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "admin@test.com' or '1'='1") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "'--") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "' #") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "'/*") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "' or true --") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if strings.Contains(email, "\") or (\"1\"=\"1") {
		message := "email format is invalid"
		errs = append(errs, errors.New(message))
	}

	if email != strings.ToLower(email) {
		message := "email addresses cannot be in uppercase."
		errs = append(errs, errors.New(message))
	}

	disposable, err := IsDisposable(email)
	if err != nil {
		message := "Internal error"
		errs = append(errs, errors.New(message))
	}

	if disposable {
		message := "temporary emails are not allowed"
		errs = append(errs, errors.New(message))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func loadDomains() {
	domains = make(map[string]struct{})

	resp, err := http.Get("https://disposable.github.io/disposable-email-domains/domains.txt")
	if err != nil {
		loadErr = err
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		domain := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if domain != "" {
			domains[domain] = struct{}{}
		}
	}

	loadErr = scanner.Err()
}

func IsDisposable(email string) (bool, error) {
	once.Do(loadDomains)

	if loadErr != nil {
		return false, loadErr
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return true, nil
	}

	domain := strings.ToLower(strings.TrimSpace(parts[1]))

	_, exists := domains[domain]

	return exists, nil
}
