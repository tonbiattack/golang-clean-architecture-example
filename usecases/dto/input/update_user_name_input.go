// Package input はユースケースへの入力DTO（Data Transfer Object）を定義します。
// DTOはレイヤー間のデータ転送に使用されます。
package input

// UpdateUserNameInput はユーザー名更新ユースケースへの入力DTOです。
type UpdateUserNameInput struct {
	UserID  string // 更新対象のユーザーID
	NewName string // 新しいユーザー名
}
