package usecase

import (
	"context"
	"log/slog"
)

type GetStoreStatus interface {
	DoByTextPost(ctx context.Context, tpg TextPostsGetter) (storeStatus, error)
}

var _ GetStoreStatus = (*getStoreStatus)(nil)

type getStoreStatus struct {
	storeStatusChecker StoreStatusChecker
}

func NewGetStoreStatus(ssc StoreStatusChecker) GetStoreStatus {
	return &getStoreStatus{
		storeStatusChecker: ssc,
	}
}

func (g *getStoreStatus) DoByTextPost(ctx context.Context, tpg TextPostsGetter) (storeStatus, error) {
	textPosts, err := tpg.GetTextPosts(ctx)
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
	GetTextPosts(ctx context.Context) ([]TextPost, error)
}

type CheckStatusInput struct {
	Posts []string
}

type storeStatus string

const (
	StoreStatusOpen   storeStatus = "open"
	StoreStatusClosed storeStatus = "closed"
)

type StoreStatusChecker interface {
	CheckStatus(ctx context.Context, input CheckStatusInput) (storeStatus, error)
}
