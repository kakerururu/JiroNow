package usecase

import (
	"context"
	"log/slog"
)

type GetStoreStatus interface {
	Do(ctx context.Context, XAccountID string) (string, error)
}

var _ GetStoreStatus = (*getStoreStatus)(nil)

type getStoreStatus struct {
	textPostGetter     TextPostsGetter
	storeStatusChecker StoreStatusChecker
}

func NewGetStoreStatus(tpg TextPostsGetter, ssc StoreStatusChecker) GetStoreStatus {
	return &getStoreStatus{
		textPostGetter:     tpg,
		storeStatusChecker: ssc,
	}
}

func (g *getStoreStatus) Do(ctx context.Context, XAccountID string) (string, error) {
	tweets, err := g.textPostGetter.GetTextPosts(XAccountID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get text posts", "error", err)
		return "", err
	}

	posts := make([]string, len(tweets))
	for i, tweet := range tweets {
		posts[i] = tweet.Post
	}

	status, err := g.storeStatusChecker.CheckStatus(CheckStatusInput{
		Posts: posts,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to check status", "error", err)
		return "", err
	}
	return status, nil
}

type TextPost struct {
	Post string
}
type TextPostsGetter interface {
	GetTextPosts(XAccountID string) ([]TextPost, error)
}

type CheckStatusInput struct {
	Posts []string
}

type StoreStatusChecker interface {
	CheckStatus(input CheckStatusInput) (string, error)
}
