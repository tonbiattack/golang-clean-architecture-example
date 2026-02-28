// Package infrastructure はインフラストラクチャ層を提供します。
// このファイルはデータベース接続の設定と初期化を行います。
package infrastructure

import (
	"database/sql"
	"fmt"

	// MySQL ドライバーを読み込み
	_ "github.com/go-sql-driver/mysql"
)

// DB はグローバルなデータベース接続を保持します（非推奨：DIで注入されるべき）
var DB *sql.DB

// NewDB は新しいMySQLデータベース接続を確立します。
// 接続文字列を構築し、データベースへの接続をテストして返します。
func NewDB() (*sql.DB, error) {
	// データベース接続パラメータ（本番環境では環境変数から取得すべき）
	dbUser := "user"
	dbPassword := "password"
	dbName := "exampledb"
	dbHost := "127.0.0.1"
	dbPort := "3306"

	// データソース名（DSN）を構築
	// parseTime=true は時間型を自動的にtime.Timeに変換
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// MySQLデータベースを開く
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// 接続テスト：データベースへの接続が有効かどうか確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}
