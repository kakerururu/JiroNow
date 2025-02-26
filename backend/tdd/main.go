package main

type TextPost struct {
	Post string
}

type TextPostGetter interface {
	GetTextPosts() ([]TextPost, error)
}

type CheckStatusInput struct {
	posts []string
}

type StoreStatusChecker interface {
	CheckStatus(input CheckStatusInput) (bool, error)
}

type GetStoreStatus struct {
	textPostGetter     TextPostGetter
	storeStatusChecker StoreStatusChecker
}

func NewGetStoreStatus(tpg TextPostGetter, ssc StoreStatusChecker) *GetStoreStatus {
	return &GetStoreStatus{
		textPostGetter:     tpg,
		storeStatusChecker: ssc,
	}
}

// 店舗名（@）,時間,→営業状況
func (g *GetStoreStatus) Do() (bool, error) {

	tweets, err := g.textPostGetter.GetTextPosts()
	if err != nil {
		return false, err
	}
	post := make([]string, len(tweets))
	for i, tweet := range tweets {
		post[i] = tweet.Post
	}
	CheckStatusInput := CheckStatusInput{
		posts: post,
	}

	result, err := g.storeStatusChecker.CheckStatus(CheckStatusInput)
	if err != nil {
		return false, err
	}
	return result, nil
}

func main() {
	tgg := &TweetGetter{}
	grg := &GeminiResponseGetter{}
	hey := NewGetStoreStatus(tgg, grg)
	hey.Do()
}
