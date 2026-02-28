// Package models はデータベーステーブルに対応するデータモデルを定義します。
// インフラストラクチャ層のモデルは、ドメインエンティティとの間で変換を行います。
package models

import "golang-clean-architecture-example/domain/entities"

// User はデータベースのusersテーブルに対応するモデルです。
// フィールドはパブリックであり、database/sqlパッケージのスキャンに使用されます。
type User struct {
	ID   string // ユーザーID（プライマリキー）
	Name string // ユーザー名
}

// FromDomainModel はドメインエンティティをインフラモデルに変換します。
// データベースに保存する前に使用します。
func FromDomainModel(m *entities.User) *User {
	return &User{
		ID:   m.GetID(),
		Name: m.GetName(),
	}
}

// ToDomainModel はインフラモデルをドメインエンティティに変換します。
// データベースから取得した後に使用します。
func (m *User) ToDomainModel() *entities.User {
	return entities.NewUser(m.ID, m.Name)
}
