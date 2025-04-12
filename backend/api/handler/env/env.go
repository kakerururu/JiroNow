package env

import "os"

type Config struct {
	XBearerToken  string
	GeminiToken   string
	XHostURL      string
	GeminiHostURL string
	GeminiPathURL string
}

func NewEnv() *Config {
	return &Config{
		XBearerToken:  getEnv("X_BEARER_TOKEN", ""),
		GeminiToken:   getEnv("GEMINI_TOKEN", ""),
		XHostURL:      getEnv("X_HOST_URL", ""),
		GeminiHostURL: getEnv("GEMINI_HOST_URL", ""),
		GeminiPathURL: getEnv("GEMINI_PATH_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
