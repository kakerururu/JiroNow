package infra

import (
	"JiroNow/usecase"
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGeminiClient_CheckStatus(t *testing.T) {
	type fields struct {
		googleApiKey string
		baseURL      url.URL
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
		wantErr bool
	}{
		{
			name: "正常系_エラーなく営業状況を取得できる",
			fields: fields{
				googleApiKey: "test",
				baseURL:      url.URL{},
				httpClient: &http.Client{
					Transport: http.DefaultTransport,
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CheckStatusInput{
					Posts: []string{("test")},
				},
			},
			want:    usecase.StoreStatusClosed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GeminiClient{
				googleApiKey: tt.fields.googleApiKey,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
			}
			got, err := s.CheckStatus(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeminiClient.CheckStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeminiClient.CheckStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
