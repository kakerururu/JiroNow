package main

import (
	"JiroNow/oas"
	"context"
	"net/http"
)

// func main() {
// 	demo_name := "demo"
// 	tgg := infra.NewTweetGetter()
// 	grg := infra.NewGeminiResponseGetter()
// 	hey := usecase.NewGetStoreStatus(tgg, grg)
// 	hey.Do(demo_name)
// }

var _ oas.Handler = (*Handler)(nil)

// oasさんが用意したインタフェースを満たすようなハンドラを作る。そこでユースケースなどに処理を投げる

type Handler struct {
}

// NewError implements oas.Handler.
func (h *Handler) NewError(ctx context.Context, err error) *oas.ErrRespStatusCode {
	panic("unimplemented")
}

// StoresStoreNameStatusGet implements oas.Handler.
func (h *Handler) StoresStoreNameStatusGet(ctx context.Context, params oas.StoresStoreNameStatusGetParams) (*oas.StoresStoreNameStatusGetOK, error) {
	panic("unimplemented")
}

func main() {

	srv, err := oas.NewServer(
		&Handler{},
	)
	if err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		panic(err)
	}
}
