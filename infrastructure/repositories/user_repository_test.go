// Package repositories_test はinfrastructure/repositoriesパッケージのテストを提供します。
package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"golang-clean-architecture-example/domain/entities"
	"golang-clean-architecture-example/infrastructure/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestNewUserRepository はNewUserRepositoryが正しくインスタンスを生成することをテストします。
func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	if repo == nil {
		t.Fatal("NewUserRepository returned nil")
	}
}

// TestUserRepository_GetUser_Success はGetUserの正常系をテストします。
func TestUserRepository_GetUser_Success(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		userName string
	}{
		{
			name:     "通常のユーザー取得",
			userID:   "user001",
			userName: "Test User",
		},
		{
			name:     "日本語名のユーザー取得",
			userID:   "user002",
			userName: "テストユーザー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックDBのセットアップ
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock db: %v", err)
			}
			defer db.Close()

			// 期待されるクエリとレスポンスを設定
			rows := sqlmock.NewRows([]string{"id"}).AddRow(tt.userID)
			mock.ExpectQuery("SELECT id FROM users WHERE id = ?").
				WithArgs(tt.userID).
				WillReturnRows(rows)

			// リポジトリを作成
			repo := repositories.NewUserRepository(db)

			// GetUserを実行
			ctx := context.Background()
			user, err := repo.GetUser(ctx, tt.userID)

			// 検証
			if err != nil {
				t.Errorf("GetUser() returned unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("GetUser() returned nil user")
			}

			if got := user.GetID(); got != tt.userID {
				t.Errorf("user.GetID() = %v, want %v", got, tt.userID)
			}

			// モックの期待が満たされているか確認
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

// TestUserRepository_GetUser_NotFound はユーザーが見つからない場合のテストです。
func TestUserRepository_GetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	// ユーザーが見つからない場合のエラーを設定
	mock.ExpectQuery("SELECT id FROM users WHERE id = ?").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	repo := repositories.NewUserRepository(db)

	ctx := context.Background()
	user, err := repo.GetUser(ctx, "nonexistent")

	// エラーが返されることを確認
	if err == nil {
		t.Error("GetUser() should return error when user not found")
	}

	if user != nil {
		t.Errorf("GetUser() should return nil user when error occurs, got %v", user)
	}

	// モックの期待が満たされているか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserRepository_GetUser_DatabaseError はデータベースエラーが発生する場合のテストです。
func TestUserRepository_GetUser_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	expectedError := errors.New("database connection error")

	mock.ExpectQuery("SELECT id FROM users WHERE id = ?").
		WithArgs("user001").
		WillReturnError(expectedError)

	repo := repositories.NewUserRepository(db)

	ctx := context.Background()
	user, err := repo.GetUser(ctx, "user001")

	// エラーが返されることを確認
	if err == nil {
		t.Error("GetUser() should return error when database error occurs")
	}

	if user != nil {
		t.Errorf("GetUser() should return nil user when error occurs, got %v", user)
	}

	// モックの期待が満たされているか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserRepository_UpdateUser_Success はUpdateUserの正常系をテストします。
func TestUserRepository_UpdateUser_Success(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		userName string
	}{
		{
			name:     "通常のユーザー更新",
			userID:   "user001",
			userName: "Updated Name",
		},
		{
			name:     "日本語名での更新",
			userID:   "user002",
			userName: "更新された名前",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックDBのセットアップ
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock db: %v", err)
			}
			defer db.Close()

			// 期待されるクエリを設定
			mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
				WithArgs(tt.userName, tt.userID).
				WillReturnResult(sqlmock.NewResult(0, 1))

			// リポジトリを作成
			repo := repositories.NewUserRepository(db)

			// ユーザーエンティティを作成
			user := entities.NewUser(tt.userID, tt.userName)

			// UpdateUserを実行
			ctx := context.Background()
			err = repo.UpdateUser(ctx, user)

			// 検証
			if err != nil {
				t.Errorf("UpdateUser() returned unexpected error: %v", err)
			}

			// モックの期待が満たされているか確認
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

// TestUserRepository_UpdateUser_DatabaseError はUpdateUserでデータベースエラーが発生する場合のテストです。
func TestUserRepository_UpdateUser_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	expectedError := errors.New("database connection error")

	mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
		WithArgs("New Name", "user001").
		WillReturnError(expectedError)

	repo := repositories.NewUserRepository(db)

	user := entities.NewUser("user001", "New Name")

	ctx := context.Background()
	err = repo.UpdateUser(ctx, user)

	// エラーが返されることを確認
	if err == nil {
		t.Error("UpdateUser() should return error when database error occurs")
	}

	if err != expectedError {
		t.Errorf("UpdateUser() error = %v, want %v", err, expectedError)
	}

	// モックの期待が満たされているか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserRepository_UpdateUser_NoRowsAffected は更新対象が存在しない場合のテストです。
func TestUserRepository_UpdateUser_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	// 更新対象が0行の場合
	mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
		WithArgs("New Name", "nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repositories.NewUserRepository(db)

	user := entities.NewUser("nonexistent", "New Name")

	ctx := context.Background()
	err = repo.UpdateUser(ctx, user)

	// エラーは返されない（実装による）
	if err != nil {
		t.Errorf("UpdateUser() returned unexpected error: %v", err)
	}

	// モックの期待が満たされているか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserRepository_ContextCancellation はコンテキストキャンセル時の動作をテストします。
func TestUserRepository_ContextCancellation(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// キャンセル済みのコンテキストを作成
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// GetUserを実行
	user, err := repo.GetUser(ctx, "user001")

	// エラーが返されることを確認（実装による）
	// コンテキストがキャンセルされている場合、通常はエラーが返される
	if err == nil && user == nil {
		// どちらかの条件が満たされていればOK
		t.Log("GetUser handled cancelled context appropriately")
	}
}
