package main

import (
	"errors"
	"testing"
)

var _ TextPostGetter = (*MockTweetGetter)(nil)

type MockTweetGetter struct {
	tweets []string
	err    error
}

func (m *MockTweetGetter) GetTextPosts() ([]TextPost, error) {
	posts := make([]TextPost, len(m.tweets))
	for i, tweet := range m.tweets {
		posts[i] = TextPost{Post: tweet}
	}
	return posts, m.err
}

var _ StoreStatusChecker = (*MockGeminiResponseGetter)(nil)

type MockGeminiResponseGetter struct {
	result bool
	err    error
}

func (m *MockGeminiResponseGetter) CheckStatus(input CheckStatusInput) (bool, error) {
	return m.result, m.err
}

// func TestGetStoreStatus_Do(t *testing.T) {
// 	type fields struct {
// 		tg TweetGetter
// 		gr GeminiResponseGetter
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "正常系_ツイートをエラーなく取得し、AIのレスポンスもエラーなく取得でき、trueが返ってきたらtrueを返す",
// 			fields: fields{
// 				tg: &MockTweetGetter{
// 					tweets: []string{"tweet1", "tweet2"},
// 					err:    nil,
// 				},
// 				gr: &MockGeminiResponseGetter{
// 					result: true,
// 					err:    nil,
// 				},
// 			},
// 			want:    true,
// 			wantErr: false,
// 		},
// 		{
// 			name: "異常系_ツイートを取得する際にエラーが発生し、エラーを返す",
// 			fields: fields{
// 				tg: &MockTweetGetter{
// 					tweets: nil,
// 					err:    errors.New("ツイート取得エラー"),
// 				},
// 				gr: &MockGeminiResponseGetter{
// 					result: true,
// 					err:    nil,
// 				},
// 			},
// 			want:    false, //boolの初期値はfalse
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			g := &GetStoreStatus{
// 				tg: tt.fields.tg,
// 				gr: tt.fields.gr,
// 			}
// 			got, err := g.Do()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetStoreStatus.Do() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("GetStoreStatus.Do() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestGetStoreStatus_Do(t *testing.T) {
	type fields struct {
		textPostGetter     TextPostGetter
		storeStatusChecker StoreStatusChecker
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "正常系_お店の投稿を取得し、判定もエラーなく取得でき、trueが返ってきたらtrueを返す",
			fields: fields{
				textPostGetter: &MockTweetGetter{
					tweets: []string{"tweet1", "tweet2"},
					err:    nil,
				},
				storeStatusChecker: &MockGeminiResponseGetter{
					result: true,
					err:    nil,
				},
			},
			want:    true,
			wantErr: false,
		}, {
			name: "異常系_お店の投稿を取得する際にエラーが発生し、エラーを返す",
			fields: fields{
				textPostGetter: &MockTweetGetter{
					tweets: nil,
					err:    errors.New("ツイート取得エラー"),
				},
				storeStatusChecker: &MockGeminiResponseGetter{
					result: true,
					err:    nil,
				},
			},
			want:    false,
			wantErr: true,
		}, {
			name: "異常系_お店の投稿をエラー無く取得し、判定でエラーが発生し、エラーを返す",
			fields: fields{
				textPostGetter: &MockTweetGetter{
					tweets: []string{"tweet1", "tweet2"},
					err:    nil,
				},
				storeStatusChecker: &MockGeminiResponseGetter{
					result: false,
					err:    errors.New("判定エラー"),
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetStoreStatus{
				textPostGetter:     tt.fields.textPostGetter,
				storeStatusChecker: tt.fields.storeStatusChecker,
			}
			got, err := g.Do()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStoreStatus.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStoreStatus.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
