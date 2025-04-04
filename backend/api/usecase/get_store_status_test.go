package usecase_test

import (
	"JiroNow/usecase"
	mockUsecase "JiroNow/usecase/mock"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_getStoreStatus_DoByTextPost(t *testing.T) {
	type mocks struct {
		mockTextPostsGetter    *mockUsecase.MockTextPostsGetter
		mockStoreStatusChecker *mockUsecase.MockStoreStatusChecker
	}

	accountNotFound := fmt.Errorf("account not found")
	notJudgeError := fmt.Errorf("not judge error")

	tests := []struct {
		name    string
		want    usecase.StoreStatus
		wantErr error
		mocks   mocks
	}{
		{
			name:    "正常系_投稿を取得し、営業状況をopenで返す",
			want:    usecase.StoreStatusOpen,
			wantErr: nil,
			mocks: func() mocks {
				ctrl := gomock.NewController(t)

				mockTextPostsGetter := mockUsecase.NewMockTextPostsGetter(ctrl)
				mockTextPostsGetter.EXPECT().GetTextPosts(t.Context()).Return([]usecase.TextPost{
					{Post: "本日は予定通りの営業です"},
				}, nil)

				mockStoreStatusChecker := mockUsecase.NewMockStoreStatusChecker(ctrl)
				mockStoreStatusChecker.EXPECT().CheckStatus(t.Context(), usecase.CheckStatusInput{
					Posts: []string{"本日は予定通りの営業です"},
				}).Return(usecase.StoreStatusOpen, nil)

				return mocks{
					mockTextPostsGetter:    mockTextPostsGetter,
					mockStoreStatusChecker: mockStoreStatusChecker,
				}
			}(),
		}, {
			name:    "異常系_TextPostGetterがエラーを返したとき、そのエラーを返す",
			want:    usecase.StoreStatusClosed, //デフォルト値をcloseとする
			wantErr: accountNotFound,
			mocks: func() mocks {
				ctrl := gomock.NewController(t)

				mockTextPostsGetter := mockUsecase.NewMockTextPostsGetter(ctrl)
				mockTextPostsGetter.EXPECT().GetTextPosts(t.Context()).Return([]usecase.TextPost{
					{Post: "本日は予定通りの営業です"},
				}, accountNotFound)

				return mocks{
					mockTextPostsGetter:    mockTextPostsGetter,
					mockStoreStatusChecker: nil,
				}
			}(),
		}, {
			name:    "異常系_投稿を取得し、StoreStatusCheckerがエラーを返したとき、そのエラーを返す",
			want:    usecase.StoreStatusClosed,
			wantErr: notJudgeError,
			mocks: func() mocks {
				ctrl := gomock.NewController(t)

				mockTextPostsGetter := mockUsecase.NewMockTextPostsGetter(ctrl)
				mockTextPostsGetter.EXPECT().GetTextPosts(t.Context()).Return([]usecase.TextPost{
					{Post: "本日は予定通りの営業です"},
				}, nil)

				mockStoreStatusChecker := mockUsecase.NewMockStoreStatusChecker(ctrl)
				mockStoreStatusChecker.EXPECT().CheckStatus(t.Context(), usecase.CheckStatusInput{
					Posts: []string{"本日は予定通りの営業です"},
				}).Return(usecase.StoreStatusClosed, notJudgeError)

				return mocks{
					mockTextPostsGetter:    mockTextPostsGetter,
					mockStoreStatusChecker: mockStoreStatusChecker,
				}
			}(),
		}, {
			name:    "異常系_投稿をエラー無く取得できたが、投稿が0件だったときエラーを返す",
			want:    usecase.StoreStatusClosed,
			wantErr: usecase.PostNotFoundError,
			mocks: func() mocks {
				ctrl := gomock.NewController(t)
				mockTextPostsGetter := mockUsecase.NewMockTextPostsGetter(ctrl)
				mockTextPostsGetter.EXPECT().GetTextPosts(t.Context()).Return([]usecase.TextPost{}, nil)
				return mocks{
					mockTextPostsGetter:    mockTextPostsGetter,
					mockStoreStatusChecker: nil,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := usecase.NewGetStoreStatus(tt.mocks.mockStoreStatusChecker)
			got, err := g.DoByTextPost(t.Context(), tt.mocks.mockTextPostsGetter)

			if err != tt.wantErr {
				t.Errorf("getStoreStatus.DoByTextPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStoreStatus.DoByTextPost() = %v, want %v", got, tt.want)
			}
		})
	}
}
