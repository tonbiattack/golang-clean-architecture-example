// Package repositories はドメインオブジェクトを永続化するためのリポジトリインターフェースを定義します。
// ドメイン層はインターフェースのみを定義し、実装はインフラストラクチャ層が行います。
// これにより依存性逆転の原則（DIP）を実現します。
package repositories

import (
	"context"
	"golang-clean-architecture-example/domain/entities"
)

// IUserRepository はUserエンティティの永続化操作を定義するインターフェースです。
// 実装はインフラストラクチャ層が提供します。
type IUserRepository interface {
	// GetUser は指定されたIDのユーザーを取得します。
	GetUser(ctx context.Context, id string) (*entities.User, error)

	// UpdateUser はユーザー情報を更新します。
	UpdateUser(ctx context.Context, user *entities.User) error
}
