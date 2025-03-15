package usecase

import (
	"context"
	"log/slog"
)

var _ GetStoreStatus = (*getStoreStatus)(nil)

type GetStoreStatus interface {
	Do(ctx context.Context, XAccountID string) (StoreStatus, error)
}

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

func (g *getStoreStatus) Do(ctx context.Context, XAccountID string) (StoreStatus, error) {
	textPosts, err := g.textPostGetter.GetTextPosts(ctx, XAccountID)

	if err != nil {

		slog.Error("failed to get text posts")
		return "", err
	}

	posts := make([]string, len(textPosts))
	for i, tweet := range textPosts {
		posts[i] = tweet.Post
	}

	status, err := g.storeStatusChecker.CheckStatus(ctx, CheckStatusInput{
		Posts: posts,
	})
	if err != nil {
		slog.Error("failed to check status")
		return "", err
	}
	return status, nil
}

type TextPost struct {
	Post string
}
type TextPostsGetter interface {
	GetTextPosts(ctx context.Context, XAccountID string) ([]TextPost, error)
}

type CheckStatusInput struct {
	Posts []string
}

type StoreStatus string

const (
	StoreStatusOpen   StoreStatus = "open"
	StoreStatusClosed StoreStatus = "closed"
)

type StoreStatusChecker interface {
	CheckStatus(ctx context.Context, input CheckStatusInput) (StoreStatus, error)
}
