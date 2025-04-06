package x

import (
	"JiroNow/usecase"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type XClient struct {
	bearerToken string
	baseURL     url.URL
	httpClient  *http.Client
}

func NewXClient(bearerToken string, baseURL url.URL, hc *http.Client) *XClient {
	return &XClient{
		bearerToken: bearerToken,
		baseURL:     baseURL,
		httpClient:  hc,
	}
}

type twitterRecentResponse struct {
	Data []struct {
		Text string `json:"text"`
	} `json:"data"`
}

var Non200Error = fmt.Errorf("non-200 response from X API")

func (x *XClient) GetXPosts(ctx context.Context, accountID string) ([]usecase.TextPost, error) {
	//URLの構築
	endpoint := x.baseURL
	endpoint.Path = "tweets/search/recent"
	query := endpoint.Query()
	query.Set("query", "from:"+accountID)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Twitter API returned non-200", "status", resp.StatusCode, "body", string(body))
		return nil, Non200Error
	}
	defer resp.Body.Close()

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

	var posts []usecase.TextPost
	for _, d := range jsonRes.Data {
		posts = append(posts, usecase.TextPost{Post: d.Text})
	}

	return posts, nil
}
