package handler

import (
	"JiroNow/handler/internal/env"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"context"
	"net/http"
)

var _ oas.Handler = (*handler)(nil)

type handler struct {
	xClient     *x.XClient
	GeminiToken string
}

func NewHandler() (*handler, error) {
	envs, err := env.NewEnv()
	if err != nil {
		return nil, err
	}
	return &handler{
		xClient:     x.NewXClient(envs.XBearerToken, http.DefaultClient),
		GeminiToken: envs.GeminiToken,
	}, nil
}

func (*handler) NewError(ctx context.Context, err error) *oas.InternalServerErrorResponseStatusCode {
	return &oas.InternalServerErrorResponseStatusCode{
		StatusCode: 500,
		Response: oas.InternalServerErrorResponse{
			Message: "internal server error",
		},
	}
}
