package handler

import (
	"JiroNow/handler/env"
	"JiroNow/infra"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"JiroNow/usecase"
	"context"
	"net/http"
	"net/url"
)

func (h *handler) StoreStatusGet(ctx context.Context, params oas.StoreStatusGetParams) (oas.StoreStatusGetRes, error) {
	var config = env.NewEnv()
	xBearerToken := config.XBearerToken
	xBaseURL := url.URL{
		Scheme: "https",
		Host:   config.XHostURL,
	}
	xClient := x.NewXClient(xBearerToken, xBaseURL, http.DefaultClient)
	xPostGetter := infra.NewXPostGetter(params.XAccountId, *xClient)

	geminiToken := config.GeminiToken
	geminiBaseURL := url.URL{
		Scheme: "https",
		Host:   config.GeminiHostURL,
		Path:   config.GeminiPathURL,
	}

	storeStatusChecker := infra.NewGeminiClientRetry(infra.NewgeminiClient(
		geminiToken,
		geminiBaseURL,
		http.DefaultClient,
	))

	uc := usecase.NewGetStoreStatus(storeStatusChecker)
	status, err := uc.DoByTextPost(ctx, xPostGetter)
	// ここにyamlの定義に合うようなエラーハンドリングを追加したいけど一旦鴨野に相談

	switch {
	case err == infra.GeminiNotJudgeError:
		return &oas.UnprocessableStoreStatusResponse{
			Message: oas.OptString{
				Value: "can't judge",
				Set:   true,
			},
		}, nil

	case err != nil:
		return nil, h.NewError(ctx, err)
	}

	switch status {
	case usecase.StoreStatusOpen:
		return &oas.GetStoreStatusResponse{
			Status: oas.NewOptGetStoreStatusResponseStatus(oas.GetStoreStatusResponseStatusOpen),
		}, nil
	case usecase.StoreStatusClosed:
		return &oas.GetStoreStatusResponse{
			Status: oas.NewOptGetStoreStatusResponseStatus(oas.GetStoreStatusResponseStatusClosed),
		}, nil
	default:
		return nil, h.NewError(ctx, err)
	}
}
