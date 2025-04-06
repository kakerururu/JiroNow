package gemini

import (
	"net/http"
	"net/url"
)

type GeminiClient struct {
	googleApiKey string
	baseURL      url.URL
	httpClient   *http.Client
}

func NewGeminiClient(key string, baseURL url.URL, hc *http.Client) *GeminiClient {
	return &GeminiClient{
		googleApiKey: key,
		baseURL:      baseURL,
		httpClient:   hc,
	}
}
