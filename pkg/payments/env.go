package payments

import "os"

// ConfigFromEnv resolves the payments base URL from the environment, mirroring
// the auth SDK convention:
//
//	PAYMENTS_ENV = "prod" | "production"  -> uses PAYMENTS_URL_PROD
//	PAYMENTS_ENV = "dev"  | "" (default)  -> uses PAYMENTS_URL_DEV
//
// If the per-env URL is unset, falls back to PAYMENTS_URL.
func ConfigFromEnv() string {
	if isProd() {
		if u := os.Getenv("PAYMENTS_URL_PROD"); u != "" {
			return u
		}
	} else {
		if u := os.Getenv("PAYMENTS_URL_DEV"); u != "" {
			return u
		}
	}
	return os.Getenv("PAYMENTS_URL")
}

func isProd() bool {
	switch os.Getenv("PAYMENTS_ENV") {
	case "prod", "production":
		return true
	default:
		return false
	}
}
