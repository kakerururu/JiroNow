package infra

import (
	"JiroNow/usecase"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var _ usecase.StoreStatusChecker = (*geminiClient)(nil)

type geminiClient struct {
	googleApiKey string
	baseURL      url.URL
	httpClient   *http.Client
}

var Non200Error = fmt.Errorf("non-200 response from API")

func NewgeminiClient(key string, baseURL url.URL, hc *http.Client) *geminiClient {
	return &geminiClient{
		googleApiKey: key,
		baseURL:      baseURL,
		httpClient:   hc,
	}
}

func (s *geminiClient) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {

	endpoint := s.baseURL
	endpoint.Path = "check_status"
	query := endpoint.Query()
	query.Set("key", s.googleApiKey)
	endpoint.RawQuery = query.Encode()

	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "以下の投稿を基に営業状況を判断してください: " + strings.Join(input.Posts, ", ")},
				},
			},
		},
	}

	bodyData, err := json.Marshal(prompt)
	if err != nil {
		return usecase.StoreStatusClosed, fmt.Errorf("failed to marshal prompt: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(bodyData))
	if err != nil {
		return usecase.StoreStatusClosed, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return usecase.StoreStatusClosed, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return usecase.StoreStatusClosed, Non200Error
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return usecase.StoreStatusClosed, err
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return usecase.StoreStatusClosed, err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return usecase.StoreStatusClosed, fmt.Errorf("invalid response structure")
	}

	statusText := response.Candidates[0].Content.Parts[0].Text
	if statusText == "open" {
		return usecase.StoreStatusOpen, nil
	}
	return usecase.StoreStatusClosed, nil
}
