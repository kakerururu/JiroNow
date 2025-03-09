package infra

import "JiroNow/usecase"

var _ usecase.TextPostGetter = (*tweetGetter)(nil)

type tweetGetter struct{}

func NewTweetGetter() *tweetGetter {
	return &tweetGetter{}
}

func (t *tweetGetter) GetTextPosts(storeName string) ([]usecase.TextPost, error) {
	// TODO: storeName を使用してテキスト投稿を取得する
	return []usecase.TextPost{{Post: "テスト投稿"}}, nil
}
