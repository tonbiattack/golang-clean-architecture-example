# Golang Clean Architecture Example

このプロジェクトは、Go言語を使用したクリーンアーキテクチャの実装例です。
Robert C. Martin（Uncle Bob）が提唱するクリーンアーキテクチャの原則に従い、保守性と拡張性が高い設計を実現しています。

## 🏗️ アーキテクチャ概要

このプロジェクトは、以下の4つの主要レイヤーで構成されています：

```
┌─────────────────────────────────────────────────┐
│          Controllers & Presenters               │  ← 外部インターフェース層
│  (HTTPリクエスト/レスポンスの処理)                   │
├─────────────────────────────────────────────────┤
│             Use Cases                           │  ← アプリケーション層
│  (ビジネスロジックの制御)                           │
├─────────────────────────────────────────────────┤
│          Domain (Entities)                      │  ← エンタープライズビジネスルール層
│  (ビジネスルールの中核)                             │
├─────────────────────────────────────────────────┤
│         Infrastructure                          │  ← フレームワーク/ドライバー層
│  (DB、外部サービスとの接続)                         │
└─────────────────────────────────────────────────┘
```

### 📁 ディレクトリ構造

```
.
├── main.go                          # アプリケーションのエントリーポイント
├── go.mod                           # Go モジュール定義
├── Makefile                         # ビルドとコード生成のコマンド
│
├── api/                             # OpenAPI仕様とコード生成
│   ├── openapi.yaml                 # API仕様
│   ├── server.go                    # 自動生成されたサーバーコード
│   └── types.go                     # 自動生成された型定義
│
├── controllers/                     # コントローラー層（外部からの入力を制御）
│   └── user_controller.go           # ユーザー関連のHTTPリクエスト処理
│
├── presenters/                      # プレゼンター層（出力を整形）
│   ├── user_presenter.go            # ユーザーレスポンスの整形
│   └── error_presenter.go           # エラーレスポンスの整形
│
├── usecases/                        # ユースケース層（アプリケーションビジネスロジック）
│   ├── update_user_name_interactor.go
│   └── dto/                         # データ転送オブジェクト
│       ├── input/
│       │   └── update_user_name_input.go
│       └── output/
│           └── update_user_name_output.go
│
├── domain/                          # ドメイン層（エンタープライズビジネスルール）
│   ├── entities/                    # ドメインエンティティ
│   │   └── user.go
│   └── repositories/                # リポジトリインターフェース
│       └── user_repository_interface.go
│
└── infrastructure/                  # インフラストラクチャ層（技術的詳細）
    ├── router.go                    # ルーティング設定
    ├── container.go                 # DI（依存性注入）コンテナ
    ├── db.go                        # データベース接続
    ├── models/                      # データベースモデル
    │   └── user_model.go
    └── repositories/                # リポジトリ実装
        └── user_repository.go
```

## 🎯 クリーンアーキテクチャの特徴

### 依存性の方向

依存性は常に外側から内側へ向かいます：

```
Infrastructure → Use Cases → Domain
Controllers    → Use Cases → Domain
Presenters     → Use Cases
```

- **Domain層**：他のどの層にも依存せず、ビジネスルールの中核を定義
- **Use Cases層**：Domainに依存し、アプリケーション固有のビジネスロジックを実装
- **Controllers/Presenters層**：Use Casesに依存し、外部とのインターフェースを提供
- **Infrastructure層**：Domainのインターフェースを実装（依存性逆転の原則）

### 依存性逆転の原則（DIP）

Domain層でリポジトリのインターフェースを定義し、Infrastructure層で実装することで、
ビジネスロジックが具体的なデータベース実装に依存しない設計を実現しています。

## 🚀 セットアップ

### 前提条件

- Go 1.22.4 以上
- MySQL 8.0 以上（またはMariaDB）
- make（オプション）

### インストール

1. リポジトリをクローン

```bash
git clone https://github.com/k-takeuchi220/golang-clean-architecture-example.git
cd golang-clean-architecture-example
```

2. 依存パッケージをインストール

```bash
go mod download
```

3. データベースの準備

```sql
CREATE DATABASE exampledb;
USE exampledb;

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO users (id, name) VALUES ('user001', 'Sample User');
```

4. データベース接続情報の設定（必要に応じて `infrastructure/db.go` を編集）

```go
dbUser := "user"
dbPassword := "password"
dbName := "exampledb"
dbHost := "127.0.0.1"
dbPort := "3306"
```

## 🏃 実行方法

### サーバーの起動

```bash
# makeを使用する場合
make run

# または直接実行
go run main.go
```

サーバーは `http://localhost:8080` で起動します。

### API呼び出し例

ユーザー名を更新する：

```bash
curl -X PUT http://localhost:8080/users/user001/update_name \
  -H "Content-Type: application/json" \
  -d '{"name": "New User Name"}'
```

レスポンス例：

```json
{
  "id": "user001",
  "name": "New User Name"
}
```

## 🧪 テスト実行

```bash
# 全テストを実行
go test ./...

# カバレッジ付きでテスト実行
go test -cover ./...

# 詳細な出力でテスト実行
go test -v ./...
```

## 🛠️ 開発

### OpenAPIからコード生成

APIの型やサーバーインターフェースを再生成する場合：

```bash
make generate-oapi
```

このコマンドは以下を実行します：
1. `oapi-codegen` ツールのインストール（未インストールの場合）
2. `api/openapi.yaml` から `api/types.go` を生成
3. `api/openapi.yaml` から `api/server.go` を生成

## 📚 使用技術

- **Go 1.22.4**: プログラミング言語
- **Echo v4**: Webフレームワーク
- **MySQL**: データベース
- **uber/dig**: 依存性注入コンテナ
- **oapi-codegen**: OpenAPIからのコード生成

## 🎓 学習ポイント

このプロジェクトから学べる主要な設計パターンと原則：

1. **クリーンアーキテクチャ**: レイヤー分離とビジネスロジックの保護
2. **依存性逆転の原則（DIP）**: インターフェースを使用した疎結合
3. **単一責任の原則（SRP）**: 各レイヤーが明確な責任を持つ
4. **依存性注入（DI）**: uber/digを使用したコンポーネント管理
5. **DTOパターン**: レイヤー間のデータ転送
6. **リポジトリパターン**: データアクセスの抽象化
7. **OpenAPI**: APIファーストな設計

## 📝 ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 👥 著者

k-takeuchi220

## 🤝 コントリビューション

プルリクエストを歓迎します！大きな変更の場合は、まずissueを開いて変更内容を議論してください。
