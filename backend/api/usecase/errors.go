package usecase

var _ GetStoreStatus = (*getStoreStatus)(nil)

type AccountNotFoundError struct {
}

func (e *AccountNotFoundError) Error() string {
	return "account not found"
}
