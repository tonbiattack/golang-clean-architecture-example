// Package entities はドメインエンティティを定義します。
// エンティティはビジネスロジックの中心となるオブジェクトで、
// 外部の技術的詳細（データベースやフレームワーク）に依存しません。
package entities

// User はユーザーを表すドメインエンティティです。
// フィールドはプライベートであり、カプセル化を保つためGetter/Setterメソッドでアクセスします。
type User struct {
	id   string // ユーザーID（一意識別子）
	name string // ユーザー名
}

// NewUser は新しいUserインスタンスを生成するファクトリ関数です。
// パラメータ:
//   - id: ユーザーID
//   - name: ユーザー名
func NewUser(id, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

// GetID はユーザーIDを取得します。
func (u *User) GetID() string {
	return u.id
}

// GetName はユーザー名を取得します。
func (u *User) GetName() string {
	return u.name
}

// SetName はユーザー名を更新します。
// パラメータ:
//   - name: 新しいユーザー名
func (u *User) SetName(name string) {
	u.name = name
}
