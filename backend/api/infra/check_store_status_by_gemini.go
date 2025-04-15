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

var _ usecase.StoreStatusChecker = (*GeminiClient)(nil)

type GeminiClient struct {
	googleApiKey string
	endpoint     url.URL
	httpClient   *http.Client
	prompt       string
}

var Non200Error = fmt.Errorf("non-200 response from API")
var GeminiNotJudgeError = fmt.Errorf("Gemini is not able to judge the status")

func NewGeminiClient(key string, hc *http.Client) *GeminiClient {
	endpoint := url.URL{
		Scheme: "https",
		Host:   "generativelanguage.googleapis.com",
		Path:   "v1beta/models/gemini-2.0-flash:generateContent",
	}
	prompt := `
		あなたは、飲食店の営業状況を判断するAIです。
		ある店舗のSNSにおける投稿を元に、その店舗の営業状況を判断してください。
		営業状況は「open」または「closed」のいずれかで答えてください。
		判断に迷った場合は、確率が高い方を選んでください。

		出力は以下の形式に従うこと。

		{
			"type": "string",
			"enum": ["open", "closed"]
		}

		open, closedのいずれかを選択し、そのテキストだけを出力してください。
		以下に出力例を示しますので、参考にしてください。

		---

		良い例:
		open

		悪い例1:
		"open"

		悪い例2:
		{
			"open"
		}

		悪い例3:
		openの確率が高いです

		悪い例4:
		open
		提示された投稿を見るに、営業状況は「open」であると考えられます。

		悪い例5:
		営業状況は「open」であると考えられます。

		---

		それでは、以下の投稿を基に営業状況を判断してください。

		---
	`
	return &GeminiClient{
		googleApiKey: key,
		endpoint:     endpoint,
		httpClient:   hc,
		prompt:       prompt,
	}
}

func (s *GeminiClient) CheckStatus(ctx context.Context, input usecase.CheckStatusInput) (usecase.StoreStatus, error) {
	endpoint := s.endpoint
	endpoint.RawQuery = "key=" + s.googleApiKey

	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": s.prompt + strings.Join(input.Posts, "\n\n")},
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
	return usecase.StoreStatusChecker(&geminiClientRetry{
		geminiClient: geminiClient,
	})
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
