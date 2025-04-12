package handler

import (
	"JiroNow/infra"
	"JiroNow/infra/x"
	"JiroNow/oas"
	"JiroNow/usecase"
	"context"
	"net/http"
	"net/url"
	"os"
)

func (h *handler) StoreStatusGet(ctx context.Context, params oas.StoreStatusGetParams) (oas.StoreStatusGetRes, error) {
	xBearerToken := os.Getenv("X_BEARER_TOKEN")
	xBaseURL := url.URL{
		Scheme: "https",
		Host:   os.Getenv("X_HOST_URL"),
	}
	xClient := x.NewXClient(xBearerToken, xBaseURL, http.DefaultClient)
	xPostGetter := infra.NewXPostGetter(params.XAccountId, *xClient)

	geminiToken := os.Getenv("GEMINI_API_KEY")
	geminiBaseURL := url.URL{
		Scheme: "https",
		Host:   os.Getenv("GEMINI_HOST_URL"),
		Path:   os.Getenv("GEMINI_PATH_URL"),
	}

	storeStatusChecker := infra.NewGeminiClientRetry(infra.NewgeminiClient(
		geminiToken,
		geminiBaseURL,
		http.DefaultClient,
	))

	uc := usecase.NewGetStoreStatus(storeStatusChecker)
	status, err := uc.DoByTextPost(ctx, xPostGetter)
	// ここにyamlの定義に合うようなエラーハンドリングを追加したいけど一旦鴨野に相談

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
