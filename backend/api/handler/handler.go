package handler

import (
	"JiroNow/oas"
	"context"
)

var _ oas.Handler = (*handler)(nil)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (*handler) NewError(ctx context.Context, err error) *oas.InternalServerErrorResponseStatusCode {
	return &oas.InternalServerErrorResponseStatusCode{
		StatusCode: 500,
		Response: oas.InternalServerErrorResponse{
			Message: oas.OptString{
				Value: err.Error(),
				Set:   false,
			},
		},
	}
}
