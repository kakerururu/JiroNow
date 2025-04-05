package x

import (
	"JiroNow/usecase"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestXClient_GetXPosts(t *testing.T) {
	type fields struct {
		bearerToken string
		baseURL     url.URL
		httpClient  *http.Client
	}
	type args struct {
		ctx       context.Context
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []usecase.TextPost
		wantErr bool
	}{
		{
			name: "正常系_ポストをエラー無く1個以上取得し、その内容を返す",
			fields: fields{
				bearerToken: "test_token",
				baseURL:     url.URL{Scheme: "https", Host: "api.twitter.com"},
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
			args: args{
				ctx:       context.Background(),
				accountID: "test_account",
			},
			want: []usecase.TextPost{
				{Post: "ポスト1"},
				{Post: "ポスト2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XClient{
				bearerToken: tt.fields.bearerToken,
				baseURL:     tt.fields.baseURL,
				httpClient:  tt.fields.httpClient,
			}
			got, err := x.GetXPosts(tt.args.ctx, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("XClient.GetXPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XClient.GetXPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
