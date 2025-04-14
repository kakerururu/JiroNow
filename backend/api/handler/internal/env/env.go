package env

import (
	"fmt"
	"log/slog"
	"os"
)

type env struct {
	XBearerToken string
	GeminiToken  string
}

func NewEnv() (*env, error) {
	xBearreToken, ok := os.LookupEnv("X_BEARER_TOKEN")
	if !ok {
		slog.Error("X_BEARER_TOKEN not set")
		return nil, fmt.Errorf("X_BEARER_TOKEN not set")
	}
	geminiToken, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		slog.Error("GEMINI_API_KEY not set")
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	return &env{
		XBearerToken: xBearreToken,
		GeminiToken:  geminiToken,
	}, nil
}
