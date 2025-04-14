package handler

import (
	"JiroNow/handler/env"
	"JiroNow/oas"
	"context"
)

var _ oas.Handler = (*handler)(nil)

type handler struct {
	XBearerToken string
	GeminiToken  string
}

func NewHandler() (*handler, error) {
	envs, err := env.NewEnv()
	if err != nil {
		return nil, err
	}
	return &handler{
		XBearerToken: envs.XBearerToken,
		GeminiToken:  envs.GeminiToken,
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
