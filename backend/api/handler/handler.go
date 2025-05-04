package handler

import (
	"JiroNow/infra"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"context"
)

var _ oas.Handler = (*handler)(nil)

type handler struct {
	xClient      *x.XClient
	geminiClient *infra.GeminiClient
}

func NewHandler(xClient *x.XClient, gClinet *infra.GeminiClient) *handler {
	return &handler{
		xClient:      xClient,
		geminiClient: gClinet,
	}
}

func (*handler) MockGet(ctx context.Context, params oas.MockGetParams) (*oas.MockStatusResponse, error) {
	// Always return "open" for the mock endpoint
	return &oas.MockStatusResponse{
		Status: "open",
	}, nil
}
