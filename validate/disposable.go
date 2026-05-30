package validate

import (
	"bufio"
	"net/http"
	"strings"
	"sync"
)

var (
	domains map[string]struct{}
	once    sync.Once
	loadErr error
)

func init() {
	go once.Do(loadDomains)
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

func isDisposable(mail string) (bool, error) {
	once.Do(loadDomains)

	if loadErr != nil {
		return false, loadErr
	}

	parts := strings.Split(mail, "@")
	if len(parts) != 2 {
		return false, nil
	}

	domain := strings.ToLower(strings.TrimSpace(parts[1]))

	_, exists := domains[domain]

	return exists, nil
}

func hasAccent(mail string) bool {
	for _, r := range mail {
		if r > 127 {
			return true
		}
	}
	return false
}
