package usecase

import (
	"context"
	"log/slog"
)

var _ GetStoreStatus = (*getStoreStatus)(nil)

type GetStoreStatus interface {
	Do(ctx context.Context) (StoreStatus, error)
}

type getStoreStatus struct {
	textPostGetter     TextPostsGetter
	storeStatusChecker StoreStatusChecker
	storeIDGetter      storeIDGetter
}

func NewGetStoreStatus(tpg TextPostsGetter, ssc StoreStatusChecker, sig storeIDGetter) GetStoreStatus {
	return &getStoreStatus{
		textPostGetter:     tpg,
		storeStatusChecker: ssc,
		storeIDGetter:      sig,
	}
}

func (g *getStoreStatus) Do(ctx context.Context) (StoreStatus, error) {
	StoreID, err := g.storeIDGetter.GetStoreID(ctx)
	if err != nil {
		slog.Error("failed to get store id")
		return "", err
	}

	textPosts, err := g.textPostGetter.GetTextPosts(ctx, StoreID)

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
	GetTextPosts(ctx context.Context, storeID string) ([]TextPost, error)
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

type storeIDGetter interface {
	GetStoreID(ctx context.Context) (string, error)
}
