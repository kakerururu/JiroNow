package x

import (
	"JiroNow/usecase"
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

	type fields struct {
		httpClient *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    []usecase.TextPost
		wantErr bool
	}{
		{
			name: "正常系_ポストをエラー無く1個以上取得し、その内容を返す",
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XClient{
				bearerToken: token,
				baseURL:     baseURL,
				httpClient:  tt.fields.httpClient,
			}
			got, err := x.GetXPosts(ctx, accountID)
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
