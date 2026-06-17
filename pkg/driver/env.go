package driver

import "os"

// ConfigFromEnv resolves the driver base URL from the environment, mirroring the
// auth SDK convention:
//
//	DRIVER_ENV = "prod" | "production"  -> uses DRIVER_URL_PROD
//	DRIVER_ENV = "dev"  | "" (default)  -> uses DRIVER_URL_DEV
//
// If the per-env URL is unset, falls back to DRIVER_URL.
func ConfigFromEnv() string {
	if isProd() {
		if u := os.Getenv("DRIVER_URL_PROD"); u != "" {
			return u
		}
	} else {
		if u := os.Getenv("DRIVER_URL_DEV"); u != "" {
			return u
		}
	}
	return os.Getenv("DRIVER_URL")
}

func isProd() bool {
	switch os.Getenv("DRIVER_ENV") {
	case "prod", "production":
		return true
	default:
		return false
	}
}
