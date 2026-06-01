package auth

import "os"

// ConfigFromEnv builds a Config from environment variables. It is the only
// place the SDK reads os.Getenv.
//
// Environment selection (dev/prod):
//
//	AUTH_ENV = "prod" | "production"  -> uses AUTH_URL_PROD
//	AUTH_ENV = "dev"  | "" (default)  -> uses AUTH_URL_DEV
//
// If the per-env URL is unset, falls back to AUTH_URL.
//
// Secrets (optional):
//
//	JWT_SECRET           -> AccessSecret   (enables offline validation)
//	INTERNAL_JWT_SECRET  -> InternalSecret (enables service tokens)
//	AUTH_USER_AGENT      -> UserAgent
func ConfigFromEnv() Config {
	return Config{
		BaseURL:        resolveBaseURL(),
		AccessSecret:   os.Getenv("JWT_SECRET"),
		InternalSecret: os.Getenv("INTERNAL_JWT_SECRET"),
		UserAgent:      os.Getenv("AUTH_USER_AGENT"),
	}
}

// IsProd reports whether AUTH_ENV selects the production environment.
func IsProd() bool {
	switch os.Getenv("AUTH_ENV") {
	case "prod", "production":
		return true
	default:
		return false
	}
}

func resolveBaseURL() string {
	if IsProd() {
		if u := os.Getenv("AUTH_URL_PROD"); u != "" {
			return u
		}
	} else {
		if u := os.Getenv("AUTH_URL_DEV"); u != "" {
			return u
		}
	}
	return os.Getenv("AUTH_URL")
}
