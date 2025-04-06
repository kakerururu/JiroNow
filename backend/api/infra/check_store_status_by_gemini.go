package infra

import (
	"JiroNow/infra/gemini"
	"JiroNow/usecase"
	"context"
)

var _ usecase.StoreStatusChecker = (*StoreStatusCheckerByGemini)(nil)

type StoreStatusCheckerByGemini struct {
	geminiClinet gemini.GeminiClient
}

func NewStoreStatusCheckerByGemini(geminiClinet gemini.GeminiClient) *StoreStatusCheckerByGemini {
	return &StoreStatusCheckerByGemini{
		geminiClinet: geminiClinet,
	}
}

func (s *StoreStatusCheckerByGemini) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {
	//ここに実装を書く
	return usecase.StoreStatusClosed, nil // Placeholder implementation
}
