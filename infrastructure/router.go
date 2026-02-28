// Package infrastructure はインフラストラクチャ層を提供します。
// このパッケージは、外部フレームワーク（Echo）やデータベース接続などの技術的詳細を扱います。
package infrastructure

import (
	"golang-clean-architecture-example/api"
	"golang-clean-architecture-example/controllers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server は HTTP サーバーを表す構造体です。
// OpenAPI で定義された ServerInterface を実装し、コントローラーに処理を委譲します。
type Server struct {
	userController *controllers.UserController
}

// NewServer は新しい Server インスタンスを生成します。
// 依存性注入により UserController を受け取ります。
func NewServer(
	userController *controllers.UserController,
) *Server {
	return &Server{
		userController: userController,
	}
}

// UpdateUserName は OpenAPI で定義された ServerInterface の実装です。
// ユーザー名更新リクエストをコントローラーに委譲します。
func (s *Server) UpdateUserName(ctx echo.Context, userID string) error {
	return s.userController.UpdateUserName(ctx, userID)
}

// InitRouter は Echo フレームワークを使用してHTTPサーバーを初期化し、起動します。
// 依存性注入コンテナを構築し、ミドルウェアを設定し、APIルートを登録します。
// サーバーはポート8080で起動します。
func InitRouter() {
	// Echo インスタンスを作成
	e := echo.New()

	// ロギングミドルウェアを追加（リクエストとレスポンスをログ出力）
	e.Use(middleware.Logger())
	// パニックリカバリーミドルウェアを追加（予期しないエラーから復旧）
	e.Use(middleware.Recover())

	// DIコンテナを構築
	container := BuildContainer()

	// コンテナからServerインスタンスを解決
	var server *Server
	if err := container.Invoke(func(s *Server) {
		server = s
	}); err != nil {
		log.Fatalf("Error resolving dependencies: %v", err)
	}

	// OpenAPI で生成されたハンドラーを登録
	api.RegisterHandlers(e, server)

	// サーバーを起動（ポート8080でリッスン）
	e.Logger.Fatal(e.Start(":8080"))
}
