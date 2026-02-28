// Package infrastructure はインフラストラクチャ層を提供します。
// このファイルは依存性注入（DI）コンテナの設定を行います。
package infrastructure

import (
	"golang-clean-architecture-example/controllers"
	"golang-clean-architecture-example/infrastructure/repositories"
	"golang-clean-architecture-example/presenters"
	"golang-clean-architecture-example/usecases"

	"go.uber.org/dig"
)

// BuildContainer はアプリケーション全体の依存性注入コンテナを構築します。
// uber/dig ライブラリを使用して、各レイヤーのコンポーネントを登録します。
// これにより、依存関係の自動解決と疎結合なアーキテクチャを実現します。
func BuildContainer() *dig.Container {
	// 新しいDIコンテナを作成
	container := dig.New()

	// インフラストラクチャ層のコンポーネントを登録
	container.Provide(NewServer) // HTTPサーバー
	container.Provide(NewDB)     // データベース接続

	// コントローラー層：HTTPリクエストを処理
	container.Provide(controllers.NewUserController)

	// プレゼンター層：レスポンスの整形と出力
	container.Provide(presenters.NewUserPresenter)
	container.Provide(presenters.NewErrorPresenter)

	// ユースケース層：ビジネスロジックの実行
	container.Provide(usecases.NewUpdateUserNameInteractor)

	// リポジトリ層：データの永続化
	container.Provide(repositories.NewUserRepository)

	return container
}
