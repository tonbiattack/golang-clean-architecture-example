// Package usecases はアプリケーションのビジネスロジックを定義します。
// ユースケース（インタラクター）は、アプリケーション固有の業務フローを管理します。
package usecases

import (
	"context"
	"golang-clean-architecture-example/domain/repositories"
	"golang-clean-architecture-example/usecases/dto/input"
	"golang-clean-architecture-example/usecases/dto/output"
)

// IUpdateUserNameInteractor はユーザー名更新のユースケースを定義するインターフェースです。
type IUpdateUserNameInteractor interface {
	// Execute はユーザー名更新のビジネスロジックを実行します。
	Execute(ctx context.Context, r *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error)
}

// UpdateUserNameInteractor はIUpdateUserNameInteractorインターフェースの実装です。
// ユーザー名更新のビジネスロジックをカプセル化します。
type UpdateUserNameInteractor struct {
	userRepository repositories.IUserRepository // ユーザーリポジトリ（依存性逆転）
}

// NewUpdateUserNameInteractor は新しいUpdateUserNameInteractorインスタンスを生成します。
// パラメータ:
//   - userRepository: ユーザーリポジトリのインターフェース
func NewUpdateUserNameInteractor(
	userRepository repositories.IUserRepository,
) IUpdateUserNameInteractor {
	return &UpdateUserNameInteractor{
		userRepository: userRepository,
	}
}

// Execute はユーザー名更新のユースケースを実行します。
// 処理フロー:
//  1. ユーザーをリポジトリから取得
//  2. ユーザー名を更新
//  3. 更新したユーザーをリポジトリに保存
//  4. 結果を返却
func (i *UpdateUserNameInteractor) Execute(ctx context.Context, input *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error) {
	// ユーザーを取得
	user, err := i.userRepository.GetUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// ユーザー名を更新（ドメインロジック）
	user.SetName(input.NewName)

	// 更新を永続化
	err = i.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 出力DTOを生成
	output := &output.UpdateUserNameOutput{
		User: user,
	}

	return output, nil
}
