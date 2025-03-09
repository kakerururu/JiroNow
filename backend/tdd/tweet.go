package main

import "database/sql"

var _ TextPostGetter = (*TweetGetter)(nil)

type TweetGetter struct{}

func (t *TweetGetter) GetTextPosts(storeName string) ([]TextPost, error) {
	// TODO: storeName を使用してテキスト投稿を取得する
	return []TextPost{{Post: "テスト投稿"}}, nil
}

var _ StoreNameGetter = (*StoreOfficialPageGetter)(nil)

type StoreOfficialPageGetter struct{}

func (s *StoreOfficialPageGetter) GetStoreName() (string, error) {
	return "仙台駅前店", nil
}

// データベースから店舗名を取得する場合はこう書くよな

// まずリポジトリのインタフェース
type StoreRepositorier interface {
	FindByID(id int) (string, error)
}

// リポジトリの実装
type StoreRepository struct {
	db *sql.DB
}

func (s *StoreRepository) FindByID(id int) (string, error) {
	// 本来はデータベースから取得する処理
	return "仙台駅前店", nil
}
