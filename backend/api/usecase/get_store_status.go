package usecase

import (
	"context"
	"log/slog"
)

type GetStoreStatus interface {
	DoByTextPost(ctx context.Context, tpg TextPostsGetter) (StoreStatus, error)
}

var _ GetStoreStatus = (*getStoreStatus)(nil)

type getStoreStatus struct {
	storeStatusChecker StoreStatusChecker
}

func NewGetStoreStatus(ssc StoreStatusChecker) *getStoreStatus {
	return &getStoreStatus{
		storeStatusChecker: ssc,
	}
}

func (g *getStoreStatus) DoByTextPost(ctx context.Context, tpg TextPostsGetter) (StoreStatus, error) {
	textPosts, err := tpg.GetTextPosts(ctx)
	if err != nil {
		slog.Error("failed to get text posts")
		return StoreStatusClosed, err
	}
	}

	posts := make([]string, len(textPosts))
	for i, textPost := range textPosts {
		posts[i] = textPost.Post
	}

	status, err := g.storeStatusChecker.CheckStatus(ctx, CheckStatusInput{
		Posts: posts,
	})
	if err != nil {
		slog.Error("failed to check status")
		return StoreStatusClosed, err
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

type StoreStatus string

const (
	StoreStatusOpen   StoreStatus = "open"
	StoreStatusClosed StoreStatus = "closed"
)

type StoreStatusChecker interface {
	CheckStatus(ctx context.Context, input CheckStatusInput) (StoreStatus, error)
}
