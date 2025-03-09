package main

import "testing"

var _ TextPostGetter = (*MockTweetGetter)(nil)

type MockTweetGetter struct {
	tweets []TextPost
	err    error
}

func (m *MockTweetGetter) GetTextPosts(storeName string) ([]TextPost, error) {
	posts := make([]TextPost, len(m.tweets))
	for i, tweet := range m.tweets {
		posts[i] = tweet
	}
	return posts, m.err
}

var _ StoreNameGetter = (*MockStoreOfficialPageGetter)(nil)

type MockStoreOfficialPageGetter struct {
	storeName string
	err       error
}

func (m *MockStoreOfficialPageGetter) GetStoreName() (string, error) {
	return m.storeName, m.err
}

var _ StoreStatusChecker = (*MockStoreStatusChecker)(nil)

type MockStoreStatusChecker struct {
	result bool
	err    error
}

func (m *MockStoreStatusChecker) CheckStatus(input CheckStatusInput) (bool, error) {
	return m.result, m.err
}

func TestGetStoreStatus_Do(t *testing.T) {
	type fields struct {
		textPostGetter     TextPostGetter
		storeStatusChecker StoreStatusChecker
		storeNameGetter    StoreNameGetter
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "正常系_エラー無く、店舗名を取得し、営業情報を取得し、営業中である",
			fields: fields{
				textPostGetter:     &MockTweetGetter{tweets: []TextPost{{Post: "テスト投稿"}}, err: nil},
				storeStatusChecker: &MockStoreStatusChecker{result: true, err: nil},
				storeNameGetter:    &MockStoreOfficialPageGetter{storeName: "仙台駅前店", err: nil},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetStoreStatus{
				textPostGetter:     tt.fields.textPostGetter,
				storeStatusChecker: tt.fields.storeStatusChecker,
				storeNameGetter:    tt.fields.storeNameGetter,
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
