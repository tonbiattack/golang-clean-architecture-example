// Package presenters はレスポンスの整形と出力を負当します。
// プレゼンターは、ユースケースの出力DTOをHTTPレスポンスに変換します。
package presenters

import (
	"golang-clean-architecture-example/api"
	"golang-clean-architecture-example/usecases/dto/output"
	"net/http"

	"github.com/labstack/echo/v4"
)

// IUserPresenter はユーザー関連のレスポンスを返すインターフェースです。
type IUserPresenter interface {
	// PresentUpdateUserName はユーザー名更新の成功レスポンスを返します。
	PresentUpdateUserName(c echo.Context, output *output.UpdateUserNameOutput) error
}

// UserPresenter はIUserPresenterインターフェースの実装です。
type UserPresenter struct{}

// NewUserPresenter は新しいUserPresenterインスタンスを生成します。
func NewUserPresenter() IUserPresenter {
	return &UserPresenter{}
}

// PresentUpdateUserName はユーザー名更新の成功レスポンスをJSON形式で返します。
// ユースケースの出力DTOをOpenAPIで定義されたAPIレスポンスに変換します。
func (p *UserPresenter) PresentUpdateUserName(c echo.Context, output *output.UpdateUserNameOutput) error {
	// ユースケースの出力DTOをAPIレスポンスに変換
	response := api.UpdateUserNameResponse{
		ID:   output.User.GetID(),
		Name: output.User.GetName(),
	}

	// HTTP 200 OK でJSONレスポンスを返す
	return c.JSON(http.StatusOK, response)
}
