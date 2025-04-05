package x

import (
	"JiroNow/usecase"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type XClient interface {
	GetXPosts(ctx context.Context, accountID string) ([]usecase.TextPost, error)
}

var _ XClient = (*XClientImpl)(nil)

type XClientImpl struct {
	bearerToken string
	baseURL     string
	httpClient  *http.Client
}

func NewXClientImpl(bt string, bu string, hc *http.Client) *XClientImpl {
	return &XClientImpl{
		bearerToken: bt,
		baseURL:     bu,
		httpClient:  hc,
	}
}

func (x *XClientImpl) GetXPosts(ctx context.Context, accountID string) ([]usecase.TextPost, error) {
	endpoint := fmt.Sprintf("%s/tweets/search/recent?query=from:%s", x.baseURL, accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		slog.Error("failed to create request", "err", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+x.bearerToken)

	resp, err := x.httpClient.Do(req)
	if err != nil {
		slog.Error("failed to send request", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Twitter API returned non-200", "status", resp.StatusCode, "body", string(body))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response body", "err", err)
		return nil, err
	}

	var jsonRes twitterRecentResponse
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		slog.Error("failed to parse JSON", "err", err)
		return nil, err
	}

	return toTextPosts(jsonRes.Data), nil
}

func toTextPosts(data []TweetData) []usecase.TextPost {
	var posts []usecase.TextPost
	for _, d := range data {
		posts = append(posts, usecase.TextPost{Post: d.Text})
	}
	return posts
}
