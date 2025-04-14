package handler

import (
	"JiroNow/infra"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"JiroNow/usecase"
	"context"
	"net/http"
)

func (h *handler) StoreStatusGet(ctx context.Context, params oas.StoreStatusGetParams) (oas.StoreStatusGetRes, error) {
	xBearerToken := h.XBearerToken
	xClient := x.NewXClient(xBearerToken, http.DefaultClient)
	xPostGetter := infra.NewXPostGetter(params.XAccountId, *xClient)

	geminiToken := h.GeminiToken
	GeminiClient := infra.NewGeminiClientRetry(infra.NewgeminiClient(
		geminiToken,
		http.DefaultClient,
	))

	uc := usecase.NewGetStoreStatus(GeminiClient)
	status, err := uc.DoByTextPost(ctx, xPostGetter)

	switch err {
	case infra.GeminiNotJudgeError:
		return &oas.UnprocessableStoreStatusResponse{
			Message: "Gemini not judge error",
		}, nil
	case x.AccountNotFoundError:
		return &oas.AccountNotFoundResponse{
			Message: "Account not found",
		}, nil
	case nil:
	default:
		return nil, h.NewError(ctx, err)
	}

	switch status {
	case usecase.StoreStatusOpen:
		return &oas.GetStoreStatusResponse{
			Status: oas.GetStoreStatusResponseStatusOpen,
		}, nil
	case usecase.StoreStatusClosed:
		return &oas.GetStoreStatusResponse{
			Status: oas.GetStoreStatusResponseStatusClosed,
		}, nil
	default:
		return nil, h.NewError(ctx, err)
	}
}
