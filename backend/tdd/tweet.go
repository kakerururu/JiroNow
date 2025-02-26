package main

var _ TextPostGetter = (*TweetGetter)(nil)

type TweetGetter struct{}

// GetTextPosts implements TextPostGetter.
func (k *TweetGetter) GetTextPosts() ([]TextPost, error) {
	return []TextPost{
		{Post: "Kamono1"},
		{Post: "Kamono2"},
	}, nil
}
