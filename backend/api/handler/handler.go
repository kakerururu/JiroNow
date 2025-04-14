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

func (*handler) NewError(ctx context.Context, err error) *oas.InternalServerErrorResponseStatusCode {
	return &oas.InternalServerErrorResponseStatusCode{
		StatusCode: 500,
		Response: oas.InternalServerErrorResponse{
			Message: "internal server error",
		},
	}
}
