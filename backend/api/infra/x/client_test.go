package x_test

import (
	"JiroNow/infra/x"
	"JiroNow/usecase"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var _ http.RoundTripper = (*mockRoundTripper)(nil)

type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestXClient_GetXPosts(t *testing.T) {

	baseURL := url.URL{Scheme: "https", Host: "api.twitter.com"}
	token := "dummy_token"
	accountID := "test_account"
	ctx := t.Context()

	var httpClientError = fmt.Errorf("failed to create http client")

	type fields struct {
		httpClient *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    []usecase.TextPost
		wantErr error
	}{
		{
			name: "正常系_ポストをエラー無く1個以上取得出来た場合、その内容を返す",
			fields: fields{
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: http.StatusOK,
								Body: io.NopCloser(strings.NewReader(`{
                                    "data": [
                                        {"text": "ポスト1"},
                                        {"text": "ポスト2"}
                                    ]
                                }`)),
							}, nil
						},
					},
				},
			},

			want: []usecase.TextPost{
				{Post: "ポスト1"},
				{Post: "ポスト2"},
			},
			wantErr: nil,
		},
		{
			name: "異常系_HTTPクライアントがエラーを返した場合、そのエラーを返す",
			fields: fields{
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {

							return nil, httpClientError
						},
					},
				},
			},
			want:    nil,
			wantErr: httpClientError,
		},
		{
			name: "異常系_レスポンスが200以外で返ってきた場合、そのエラーを返す",
			fields: fields{
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: http.StatusInternalServerError,
								Body:       io.NopCloser(strings.NewReader(`{"invalid_json":`)),
							}, nil
						},
					},
				},
			},
			want:    nil,
			wantErr: x.Non200Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := x.NewXClient(token, baseURL, tt.fields.httpClient)
			got, err := x.GetXPosts(ctx, accountID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("XClient.GetXPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XClient.GetXPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
