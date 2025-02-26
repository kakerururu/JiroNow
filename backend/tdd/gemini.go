package main

var _ StoreStatusChecker = (*GeminiResponseGetter)(nil)

type GeminiResponseGetter struct{}

// CheckStatus implements StoreStatusChecker.
func (k *GeminiResponseGetter) CheckStatus(input CheckStatusInput) (bool, error) {
	panic("unimplemented")
}
