package usecase

type GetStoreStatus interface {
	Do(storeName string) (bool, error)
}

var _ GetStoreStatus = (*getStoreStatus)(nil)

type getStoreStatus struct {
	textPostGetter     TextPostGetter
	storeStatusChecker StoreStatusChecker
}

func NewGetStoreStatus(tpg TextPostGetter, ssc StoreStatusChecker) *getStoreStatus {
	return &getStoreStatus{
		textPostGetter:     tpg,
		storeStatusChecker: ssc,
	}
}

func (g *getStoreStatus) Do(storeName string) (bool, error) {

	tweets, err := g.textPostGetter.GetTextPosts(storeName)
	if err != nil {
		return false, err
	}
	post := make([]string, len(tweets))
	for i, tweet := range tweets {
		post[i] = tweet.Post
	}
	CheckStatusInput := CheckStatusInput{
		Posts: post,
	}

	result, err := g.storeStatusChecker.CheckStatus(CheckStatusInput)
	if err != nil {
		return false, err
	}
	return result, nil
}

type TextPost struct {
	Post string
}

type TextPostGetter interface {
	GetTextPosts(storeName string) ([]TextPost, error)
}

type CheckStatusInput struct {
	Posts []string
}

type StoreStatusChecker interface {
	CheckStatus(input CheckStatusInput) (bool, error)
}
