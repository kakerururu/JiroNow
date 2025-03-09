package infra

import "JiroNow/usecase"

var _ usecase.StoreStatusChecker = (*geminiResponseGetter)(nil)

type geminiResponseGetter struct{}

func NewGeminiResponseGetter() *geminiResponseGetter {
	return &geminiResponseGetter{}
}

func (k *geminiResponseGetter) CheckStatus(input usecase.CheckStatusInput) (bool, error) {
	panic("営業していません")
}
