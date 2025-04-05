package infra

import (
	"JiroNow/infra/x"
	"JiroNow/usecase"
	"context"
)

var _ usecase.TextPostsGetter = (*xPostGetter)(nil)

type xPostGetter struct {
	xAccountID string
	xClient    x.XClient
}

func NewXPostGetter(xID string, xc x.XClient) *xPostGetter {
	return &xPostGetter{
		xAccountID: xID,
		xClient:    xc,
	}
}

func (x *xPostGetter) GetTextPosts(ctx context.Context) ([]usecase.TextPost, error) {
	accountID := x.xAccountID
	posts, err := x.xClient.GetXPosts(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
