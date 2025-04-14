package infra

import (
	"JiroNow/usecase"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

var _ usecase.StoreStatusChecker = (*geminiClient)(nil)

type geminiClient struct {
	googleApiKey string
	endpoint     url.URL
	httpClient   *http.Client
}

var Non200Error = fmt.Errorf("non-200 response from API")
var GeminiNotJudgeError = fmt.Errorf("Gemini is not able to judge the status")

func NewgeminiClient(key string, hc *http.Client) *geminiClient {
	endpoint := url.URL{
		Scheme: "https",
		Host:   "generativelanguage.googleapis.com",
		Path:   "v1beta/models/gemini-2.0-flash:generateContent",
	}
	return &geminiClient{
		googleApiKey: key,
		endpoint:     endpoint,
		httpClient:   hc,
	}
}

func (s *geminiClient) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {
	// https: //generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=AIzaSyAisMmlx3Q9s0q4MhRbvpukNTVAa-AhOOk
	endpoint := s.endpoint
	endpoint.RawQuery = "key=" + s.googleApiKey

	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "以下の投稿を基に営業状況を判断してください: " + strings.Join(input.Posts, ", ") + "\n{ 'type': 'string', 'enum': ['open', 'closed'] }"},
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

	statusText := strings.TrimSuffix(response.Candidates[0].Content.Parts[0].Text, "\n")
	if statusText == "open" {
		return usecase.StoreStatusOpen, nil
	}
	if statusText == "closed" {
		return usecase.StoreStatusClosed, nil
	}
	slog.Error("invalid status text: " + statusText)

	return usecase.StoreStatusClosed, GeminiNotJudgeError
}

var _ usecase.StoreStatusChecker = (*geminiClientRetry)(nil)

type geminiClientRetry struct {
	geminiClient usecase.StoreStatusChecker
}

func NewGeminiClientRetry(geminiClient usecase.StoreStatusChecker) usecase.StoreStatusChecker {
	return &geminiClientRetry{
		geminiClient: geminiClient,
	}
}

func (g *geminiClientRetry) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {
	const maxRetries = 3
	var err error
	var status usecase.StoreStatus

	for i := 0; i < maxRetries; i++ {
		status, err = g.geminiClient.CheckStatus(ctx, input)
		if err == GeminiNotJudgeError {
			slog.Info("retrying due to Gemini No Brain Error")
			continue
		}
		return status, err
	}

	return usecase.StoreStatusClosed, err
}
