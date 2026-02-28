// Package output はユースケースからの出力DTO（Data Transfer Object）を定義します。
// DTOはレイヤー間のデータ転送に使用されます。
package output

import "golang-clean-architecture-example/domain/entities"

// UpdateUserNameOutput はユーザー名更新ユースケースからの出力DTOです。
type UpdateUserNameOutput struct {
	User *entities.User // 更新後のユーザーエンティティ
}
