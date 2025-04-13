package env

import (
	"log/slog"
	"os"
)

type env struct {
	XBearerToken string
	GeminiToken  string
}

func NewEnv() *env {
	xBearreToken, ok := os.LookupEnv("X_BEARER_TOKEN")
	if !ok {
		slog.Error("X_BEARER_TOKEN not set")
	}
	geminiToken, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		slog.Error("GEMINI_API_KEY not set")
	}
	return &env{
		XBearerToken: xBearreToken,
		GeminiToken:  geminiToken,
	}
}
