package main

type TextPost struct {
	Post string
}

type TextPostGetter interface {
	GetTextPosts(storeName string) ([]TextPost, error)
}

type StoreNameGetter interface {
	GetStoreName() (string, error)
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
	storeNameGetter    StoreNameGetter
}

func NewGetStoreStatus(tpg TextPostGetter, ssc StoreStatusChecker, sng StoreNameGetter) *GetStoreStatus {
	return &GetStoreStatus{
		textPostGetter:     tpg,
		storeStatusChecker: ssc,
		storeNameGetter:    sng, //鴨野に相談。別に一店舗だけなら、stringを外部から注入で良くないか？
		// 店舗名取得は完全にこの構造体の責務外じゃね？
	}
}

func (g *GetStoreStatus) Do() (bool, error) {
	storeName, err := g.storeNameGetter.GetStoreName()
	if err != nil {
		return false, err
	}

	tweets, err := g.textPostGetter.GetTextPosts(storeName)
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
	sng := &StoreOfficialPageGetter{}
	hey := NewGetStoreStatus(tgg, grg, sng)
	hey.Do()
}
