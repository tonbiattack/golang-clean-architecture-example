// Package entities_test はdomain/entitiesパッケージのテストを提供します。
package entities_test

import (
	"golang-clean-architecture-example/domain/entities"
	"testing"
)

// TestNewUser はNewUser関数が正しくUserインスタンスを生成することをテストします。
func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		userName string
	}{
		{
			name:     "通常のユーザー作成",
			id:       "user001",
			userName: "Test User",
		},
		{
			name:     "空文字を含むユーザー作成",
			id:       "",
			userName: "",
		},
		{
			name:     "日本語を含むユーザー作成",
			id:       "user002",
			userName: "テストユーザー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := entities.NewUser(tt.id, tt.userName)

			if user == nil {
				t.Fatal("NewUser returned nil")
			}

			if got := user.GetID(); got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}

			if got := user.GetName(); got != tt.userName {
				t.Errorf("GetName() = %v, want %v", got, tt.userName)
			}
		})
	}
}

// TestUser_GetID はGetIDメソッドが正しく動作することをテストします。
func TestUser_GetID(t *testing.T) {
	expectedID := "test-id-123"
	user := entities.NewUser(expectedID, "Test User")

	if got := user.GetID(); got != expectedID {
		t.Errorf("GetID() = %v, want %v", got, expectedID)
	}
}

// TestUser_GetName はGetNameメソッドが正しく動作することをテストします。
func TestUser_GetName(t *testing.T) {
	expectedName := "Test User"
	user := entities.NewUser("user001", expectedName)

	if got := user.GetName(); got != expectedName {
		t.Errorf("GetName() = %v, want %v", got, expectedName)
	}
}

// TestUser_SetName はSetNameメソッドが正しく動作することをテストします。
func TestUser_SetName(t *testing.T) {
	tests := []struct {
		name        string
		initialName string
		newName     string
	}{
		{
			name:        "通常の名前変更",
			initialName: "Old Name",
			newName:     "New Name",
		},
		{
			name:        "空文字への変更",
			initialName: "Old Name",
			newName:     "",
		},
		{
			name:        "日本語への変更",
			initialName: "Old Name",
			newName:     "新しい名前",
		},
		{
			name:        "長い名前への変更",
			initialName: "Old Name",
			newName:     "Very Long User Name With Many Characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := entities.NewUser("user001", tt.initialName)

			// 初期値の確認
			if got := user.GetName(); got != tt.initialName {
				t.Errorf("Initial GetName() = %v, want %v", got, tt.initialName)
			}

			// 名前を変更
			user.SetName(tt.newName)

			// 変更後の値を確認
			if got := user.GetName(); got != tt.newName {
				t.Errorf("After SetName(), GetName() = %v, want %v", got, tt.newName)
			}

			// IDが変更されていないことを確認
			if got := user.GetID(); got != "user001" {
				t.Errorf("GetID() should not change, got %v", got)
			}
		})
	}
}

// TestUser_SetNameMultipleTimes は複数回SetNameを呼び出した場合の動作をテストします。
func TestUser_SetNameMultipleTimes(t *testing.T) {
	user := entities.NewUser("user001", "Initial Name")

	names := []string{"First Change", "Second Change", "Third Change", "Final Name"}

	for _, name := range names {
		user.SetName(name)
		if got := user.GetName(); got != name {
			t.Errorf("After SetName(%v), GetName() = %v, want %v", name, got, name)
		}
	}

	// 最終的な名前の確認
	expectedFinalName := "Final Name"
	if got := user.GetName(); got != expectedFinalName {
		t.Errorf("Final GetName() = %v, want %v", got, expectedFinalName)
	}
}

// TestUser_ImmutableID はユーザーIDが不変であることをテストします。
func TestUser_ImmutableID(t *testing.T) {
	expectedID := "immutable-id"
	user := entities.NewUser(expectedID, "Initial Name")

	// IDは初期値
	if got := user.GetID(); got != expectedID {
		t.Errorf("Initial GetID() = %v, want %v", got, expectedID)
	}

	// 名前を複数回変更
	user.SetName("Name 1")
	user.SetName("Name 2")
	user.SetName("Name 3")

	// IDは変更されていないはず
	if got := user.GetID(); got != expectedID {
		t.Errorf("After multiple SetName calls, GetID() = %v, want %v", got, expectedID)
	}
}
