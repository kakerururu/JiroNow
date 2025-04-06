package infra

import (
	"JiroNow/usecase"
	"context"
	"net/http"
	"net/url"
)

var _ usecase.StoreStatusChecker = (*StoreStatusCheckerByGemini)(nil)

type StoreStatusCheckerByGemini struct {
	geminiClinet GeminiClient
}

func NewStoreStatusCheckerByGemini(geminiClinet GeminiClient) *StoreStatusCheckerByGemini {
	return &StoreStatusCheckerByGemini{
		geminiClinet: geminiClinet,
	}
}

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

func (s *StoreStatusCheckerByGemini) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {

	return usecase.StoreStatusClosed, nil // Placeholder implementation
}
