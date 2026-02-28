// Package presenters はレスポンスの整形と出力を負当します。
// このファイルはエラーレスポンスの処理を行います。
package presenters

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// IErrorPresenter はエラーレスポンスを返すインターフェースです。
type IErrorPresenter interface {
	// PresentBadRequest は400 Bad Requestレスポンスを返します。
	PresentBadRequest(c echo.Context, message string) error

	// PresentInternalServerError は500 Internal Server Errorレスポンスを返します。
	PresentInternalServerError(c echo.Context, err error) error
}

// ErrorPresenter はIErrorPresenterインターフェースの実装です。
type ErrorPresenter struct{}

// NewErrorPresenter は新しいErrorPresenterインスタンスを生成します。
func NewErrorPresenter() IErrorPresenter {
	return &ErrorPresenter{}
}

// PresentBadRequest は400 Bad RequestエラーレスポンスをJSON形式で返します。
// クライアントからの無効なリクエストに対して使用します。
func (p *ErrorPresenter) PresentBadRequest(c echo.Context, message string) error {
	// エラーメッセージを含むレスポンスを生成
	response := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	// HTTP 400 Bad Request でJSONレスポンスを返す
	return c.JSON(http.StatusBadRequest, response)
}

// PresentInternalServerError は500 Internal Server ErrorレスポンスをJSON形式で返します。
// サーバー側で発生した予期しないエラーに対して使用します。
func (p *ErrorPresenter) PresentInternalServerError(c echo.Context, err error) error {
	// エラーメッセージを含むレスポンスを生成
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	// HTTP 500 Internal Server Error でJSONレスポンスを返す
	return c.JSON(http.StatusInternalServerError, response)
}
