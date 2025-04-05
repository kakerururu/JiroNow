package x

type twitterRecentResponse struct {
	Data []TweetData `json:"data"`
}

type TweetData struct {
	Text string `json:"text"`
}
