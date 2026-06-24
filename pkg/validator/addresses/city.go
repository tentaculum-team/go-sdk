package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tentaculum-dev/go-sdk/internal/utils"
)

type CityConfig struct {
	MaxChars  int
	MinChars  int
	AllErrors bool
}

func DefaultCityConfig() CityConfig {
	return CityConfig{MaxChars: 100, MinChars: 0, AllErrors: true}
}

// City validates a city/town name for the global market. It accepts Unicode
// letters and combining marks (accents: "São Paulo", "Saint-Étienne"), digits
// ("100 Mile House"), spaces, hyphens, apostrophes ("O'Fallon"), periods and
// commas ("Washington, D.C."). Control characters and invalid UTF-8 are rejected.
func City(city string, cfg ...CityConfig) error {
	conf := DefaultCityConfig()
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	city = strings.TrimSpace(city)
	if city == "" {
		return errors.New("city cannot be empty.")
	}

	var errs []error
	n := utf8.RuneCountInString(city)

	if conf.MaxChars > 0 && n > conf.MaxChars {
		errs = append(errs, fmt.Errorf("city cannot exceed %d characters. (current %d)", conf.MaxChars, n))
	}
	if conf.MinChars > 0 && n < conf.MinChars {
		errs = append(errs, fmt.Errorf("city cannot be shorter than %d characters. (current %d)", conf.MinChars, n))
	}
	if utils.ContainsNullByte(city) || utils.ContainsControlChars(city) {
		errs = append(errs, errors.New("city cannot contain control characters."))
	}
	if utils.ContainsInvalidUTF8(city) {
		errs = append(errs, errors.New("city contains invalid UTF-8 characters."))
	}
	for _, r := range city {
		if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) ||
			strings.ContainsRune(" -'.,", r) {
			continue
		}
		errs = append(errs, fmt.Errorf("city contains an invalid character: '%c'.", r))
		break
	}

	if conf.AllErrors {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	} else if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
