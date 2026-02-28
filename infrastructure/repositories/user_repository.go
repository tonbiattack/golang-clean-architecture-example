// Package repositories はドメインリポジトリインターフェースの具体的な実装を提供します。
// ここではデータベースへの実際のアクセスを行います。
package repositories

import (
	"context"
	"database/sql"
	"golang-clean-architecture-example/domain/entities"
	"golang-clean-architecture-example/domain/repositories"
	"golang-clean-architecture-example/infrastructure/models"
)

// UserRepository はIUserRepositoryインターフェースの実装です。
// MySQLデータベースを使用してユーザー情報を永続化します。
type UserRepository struct {
	db *sql.DB // データベース接続
}

// NewUserRepository は新しいUserRepositoryインスタンスを生成します。
// インターフェース型で返すことで、呼び出し側が具体的な実装に依存しないようにします。
func NewUserRepository(db *sql.DB) repositories.IUserRepository {
	return &UserRepository{db: db}
}

// GetUser は指定されたIDのユーザーをデータベースから取得します。
// データベースモデルをドメインエンティティに変換して返します。
func (r *UserRepository) GetUser(ctx context.Context, id string) (*entities.User, error) {
	// SQLクエリを実行（注意：実際にnameカラムもSELECTすべき）
	row := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = ?", id)
	user := &models.User{}

	// 結果をスキャン
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	// インフラモデルをドメインエンティティに変換
	return user.ToDomainModel(), nil
}

// UpdateUser はユーザー情報をデータベースで更新します。
// ドメインエンティティをデータベースモデルに変換して更新します。
func (r *UserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	// ドメインエンティティをインフラモデルに変換
	model := models.FromDomainModel(user)

	// UPDATE SQLを実行
	_, err := r.db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", model.Name, model.ID)
	return err
}
