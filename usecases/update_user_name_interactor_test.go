// Package usecases_test はusecasesパッケージのテストを提供します。
package usecases_test

import (
	"context"
	"errors"
	"golang-clean-architecture-example/domain/entities"
	"golang-clean-architecture-example/usecases"
	"golang-clean-architecture-example/usecases/dto/input"
	"testing"
)

// mockUserRepository はテスト用のモックリポジトリです。
type mockUserRepository struct {
	getUserFunc    func(ctx context.Context, id string) (*entities.User, error)
	updateUserFunc func(ctx context.Context, user *entities.User) error
}

func (m *mockUserRepository) GetUser(ctx context.Context, id string) (*entities.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	return nil, errors.New("GetUser not implemented")
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, user)
	}
	return errors.New("UpdateUser not implemented")
}

// TestNewUpdateUserNameInteractor は NewUpdateUserNameInteractor が正しくインスタンスを生成することをテストします。
func TestNewUpdateUserNameInteractor(t *testing.T) {
	mockRepo := &mockUserRepository{}
	interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

	if interactor == nil {
		t.Fatal("NewUpdateUserNameInteractor returned nil")
	}
}

// TestUpdateUserNameInteractor_Execute_Success は正常系のテストです。
func TestUpdateUserNameInteractor_Execute_Success(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		initialName string
		newName     string
	}{
		{
			name:        "通常のユーザー名変更",
			userID:      "user001",
			initialName: "Old Name",
			newName:     "New Name",
		},
		{
			name:        "日本語のユーザー名変更",
			userID:      "user002",
			initialName: "旧名前",
			newName:     "新名前",
		},
		{
			name:        "空文字への変更",
			userID:      "user003",
			initialName: "Old Name",
			newName:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリのセットアップ
			mockRepo := &mockUserRepository{
				getUserFunc: func(ctx context.Context, id string) (*entities.User, error) {
					if id != tt.userID {
						t.Errorf("GetUser called with wrong id: got %v, want %v", id, tt.userID)
					}
					return entities.NewUser(tt.userID, tt.initialName), nil
				},
				updateUserFunc: func(ctx context.Context, user *entities.User) error {
					if user.GetID() != tt.userID {
						t.Errorf("UpdateUser called with wrong user id: got %v, want %v", user.GetID(), tt.userID)
					}
					if user.GetName() != tt.newName {
						t.Errorf("UpdateUser called with wrong user name: got %v, want %v", user.GetName(), tt.newName)
					}
					return nil
				},
			}

			// インタラクターを作成
			interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

			// 入力DTOを作成
			inputDTO := &input.UpdateUserNameInput{
				UserID:  tt.userID,
				NewName: tt.newName,
			}

			// Executeを実行
			ctx := context.Background()
			output, err := interactor.Execute(ctx, inputDTO)

			// エラーがないことを確認
			if err != nil {
				t.Fatalf("Execute() returned unexpected error: %v", err)
			}

			// 出力DTOの検証
			if output == nil {
				t.Fatal("Execute() returned nil output")
			}

			if output.User == nil {
				t.Fatal("output.User is nil")
			}

			if got := output.User.GetID(); got != tt.userID {
				t.Errorf("output.User.GetID() = %v, want %v", got, tt.userID)
			}

			if got := output.User.GetName(); got != tt.newName {
				t.Errorf("output.User.GetName() = %v, want %v", got, tt.newName)
			}
		})
	}
}

// TestUpdateUserNameInteractor_Execute_GetUserError はGetUserでエラーが発生する場合のテストです。
func TestUpdateUserNameInteractor_Execute_GetUserError(t *testing.T) {
	expectedError := errors.New("user not found")

	mockRepo := &mockUserRepository{
		getUserFunc: func(ctx context.Context, id string) (*entities.User, error) {
			return nil, expectedError
		},
	}

	interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

	inputDTO := &input.UpdateUserNameInput{
		UserID:  "nonexistent",
		NewName: "New Name",
	}

	ctx := context.Background()
	output, err := interactor.Execute(ctx, inputDTO)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("Execute() should return error when GetUser fails")
	}

	if err != expectedError {
		t.Errorf("Execute() error = %v, want %v", err, expectedError)
	}

	// 出力がnilであることを確認
	if output != nil {
		t.Errorf("Execute() output should be nil when error occurs, got %v", output)
	}
}

// TestUpdateUserNameInteractor_Execute_UpdateUserError はUpdateUserでエラーが発生する場合のテストです。
func TestUpdateUserNameInteractor_Execute_UpdateUserError(t *testing.T) {
	expectedError := errors.New("database connection error")

	mockRepo := &mockUserRepository{
		getUserFunc: func(ctx context.Context, id string) (*entities.User, error) {
			return entities.NewUser("user001", "Old Name"), nil
		},
		updateUserFunc: func(ctx context.Context, user *entities.User) error {
			return expectedError
		},
	}

	interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

	inputDTO := &input.UpdateUserNameInput{
		UserID:  "user001",
		NewName: "New Name",
	}

	ctx := context.Background()
	output, err := interactor.Execute(ctx, inputDTO)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("Execute() should return error when UpdateUser fails")
	}

	if err != expectedError {
		t.Errorf("Execute() error = %v, want %v", err, expectedError)
	}

	// 出力がnilであることを確認
	if output != nil {
		t.Errorf("Execute() output should be nil when error occurs, got %v", output)
	}
}

// TestUpdateUserNameInteractor_Execute_ContextCancellation はコンテキストキャンセル時の動作をテストします。
func TestUpdateUserNameInteractor_Execute_ContextCancellation(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserFunc: func(ctx context.Context, id string) (*entities.User, error) {
			// コンテキストがキャンセルされているかチェック
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return entities.NewUser("user001", "Old Name"), nil
			}
		},
	}

	interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

	inputDTO := &input.UpdateUserNameInput{
		UserID:  "user001",
		NewName: "New Name",
	}

	// キャンセル済みのコンテキストを作成
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // すぐにキャンセル

	output, err := interactor.Execute(ctx, inputDTO)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("Execute() should return error when context is cancelled")
	}

	// 出力がnilであることを確認
	if output != nil {
		t.Errorf("Execute() output should be nil when error occurs, got %v", output)
	}
}

// TestUpdateUserNameInteractor_Execute_CallOrder はメソッドの呼び出し順序が正しいことをテストします。
func TestUpdateUserNameInteractor_Execute_CallOrder(t *testing.T) {
	callOrder := []string{}

	mockRepo := &mockUserRepository{
		getUserFunc: func(ctx context.Context, id string) (*entities.User, error) {
			callOrder = append(callOrder, "GetUser")
			return entities.NewUser("user001", "Old Name"), nil
		},
		updateUserFunc: func(ctx context.Context, user *entities.User) error {
			callOrder = append(callOrder, "UpdateUser")
			return nil
		},
	}

	interactor := usecases.NewUpdateUserNameInteractor(mockRepo)

	inputDTO := &input.UpdateUserNameInput{
		UserID:  "user001",
		NewName: "New Name",
	}

	ctx := context.Background()
	_, err := interactor.Execute(ctx, inputDTO)

	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	// 呼び出し順序の確認
	if len(callOrder) != 2 {
		t.Fatalf("Expected 2 method calls, got %d", len(callOrder))
	}

	if callOrder[0] != "GetUser" {
		t.Errorf("First call should be GetUser, got %s", callOrder[0])
	}

	if callOrder[1] != "UpdateUser" {
		t.Errorf("Second call should be UpdateUser, got %s", callOrder[1])
	}
}
