package usecase

import (
	"context"
	"fmt"
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

var PostNotFoundError = fmt.Errorf("post not found")

func (g *getStoreStatus) DoByTextPost(ctx context.Context, tpg TextPostsGetter) (StoreStatus, error) {
	textPosts, err := tpg.GetTextPosts(ctx)
	if err != nil {
		slog.Error("failed to get text posts")
		return StoreStatusClosed, err
	}

	if len(textPosts) == 0 {
		slog.Error("post not found")
		return StoreStatusClosed, PostNotFoundError
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
