package env

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type env struct {
	XBearerToken string
	GeminiToken  string
	Port         int
}

func NewEnv() (*env, error) {
	xBearerToken, ok := os.LookupEnv("X_BEARER_TOKEN")
	if !ok {
		slog.Error("X_BEARER_TOKEN not set")
		return nil, fmt.Errorf("X_BEARER_TOKEN not set")
	}
	geminiToken, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		slog.Error("GEMINI_API_KEY not set")
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Warn("PORT not set, using default value 8080")
		port = "8080"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		slog.Error("Invalid PORT value, must be an integer between 1 and 65535")
		return nil, fmt.Errorf("invalid PORT value: %s", port)
	}
	return &env{
		XBearerToken: xBearerToken,
		GeminiToken:  geminiToken,
		Port:         portInt,
	}, nil
}
