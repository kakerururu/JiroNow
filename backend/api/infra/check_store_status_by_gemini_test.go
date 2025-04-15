package infra_test

import (
	"JiroNow/infra"
	"JiroNow/usecase"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func TestGeminiClient_CheckStatus(t *testing.T) {
	var HTTPClientError = fmt.Errorf("HTTP Client Error")

	type fields struct {
		googleApiKey string
		httpClient   *http.Client
	}
	type args struct {
		ctx   context.Context
		input usecase.CheckStatusInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.StoreStatus
		wantErr error
	}{
		{
			name: "正常系_エラーなく営業状況を取得できる",
			fields: fields{
				googleApiKey: "test",
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: http.StatusOK,
								Body: io.NopCloser(strings.NewReader(`{
                "candidates": [
                    {
                        "content": {
                            "parts": [
                                {
                                    "text": "closed"
                                }
                            ]
                        }
                    }
                ]
            }`)),
							}, nil
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CheckStatusInput{
					Posts: []string{("test")},
				},
			},
			want:    usecase.StoreStatusClosed,
			wantErr: nil,
		}, {
			name: "異常系_HTTPステータスコードが200以外の場合",
			fields: fields{
				googleApiKey: "test",
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: http.StatusInternalServerError,
								Body:       io.NopCloser(strings.NewReader("")),
							}, nil
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CheckStatusInput{
					Posts: []string{"test"},
				},
			},
			want:    usecase.StoreStatusClosed,
			wantErr: infra.Non200Error,
		},

		{
			name: "異常系_HTTPクライアントがエラーを返した場合",
			fields: fields{
				googleApiKey: "test",
				httpClient: &http.Client{
					Transport: &mockRoundTripper{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							return nil, HTTPClientError
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CheckStatusInput{
					Posts: []string{"test"},
				},
			},
			want:    usecase.StoreStatusClosed,
			wantErr: HTTPClientError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := infra.NewGeminiClient(tt.fields.googleApiKey, tt.fields.httpClient)
			got, err := s.CheckStatus(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GeminiClient.CheckStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeminiClient.CheckStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
