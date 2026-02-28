// Package controllers はHTTPリクエストの制御を行います。
// コントローラーは、リクエストの解析、ユースケースの呼び出し、
// プレゼンターへの委譲を行います。
package controllers

import (
	"golang-clean-architecture-example/api"
	"golang-clean-architecture-example/presenters"
	"golang-clean-architecture-example/usecases"
	"golang-clean-architecture-example/usecases/dto/input"

	"github.com/labstack/echo/v4"
)

// UserController はユーザー関連のHTTPリクエストを制御します。
type UserController struct {
	updateUserNameInteractor usecases.IUpdateUserNameInteractor // ユーザー名更新ユースケース
	userPresenter            presenters.IUserPresenter          // ユーザーレスポンスのプレゼンター
	errorPresenter           presenters.IErrorPresenter         // エラーレスポンスのプレゼンター
}

// NewUserController は新しいUserControllerインスタンスを生成します。
// パラメータ:
//   - updateUserNameInteractor: ユーザー名更新ユースケース
//   - userPresenter: ユーザーレスポンスのプレゼンター
//   - errorPresenter: エラーレスポンスのプレゼンター
func NewUserController(
	updateUserNameInteractor usecases.IUpdateUserNameInteractor,
	userPresenter presenters.IUserPresenter,
	errorPresenter presenters.IErrorPresenter,
) *UserController {
	return &UserController{
		updateUserNameInteractor: updateUserNameInteractor,
		userPresenter:            userPresenter,
		errorPresenter:           errorPresenter,
	}
}

// UpdateUserName はユーザー名更新のHTTPリクエストを処理します。
// 処理フロー:
//  1. リクエストボディをパース
//  2. 入力DTOを生成
//  3. ユースケースを実行
//  4. プレゼンターにレスポンスを委譲
func (c *UserController) UpdateUserName(e echo.Context, userID string) error {
	// リクエストボディをパース
	req := api.UpdateUserNameRequest{}
	if err := e.Bind(&req); err != nil {
		// パースエラー時はBadRequestを返す
		return c.errorPresenter.PresentBadRequest(e, "invalid request")
	}

	// 入力DTOを生成
	input := &input.UpdateUserNameInput{
		UserID:  userID,
		NewName: req.Name,
	}
	ctx := e.Request().Context()

	// ユースケースを実行
	output, err := c.updateUserNameInteractor.Execute(ctx, input)
	if err != nil {
		// エラー時はInternalServerErrorを返す
		return c.errorPresenter.PresentInternalServerError(e, err)
	}

	// 成功時はユーザープレゼンターでレスポンスを返す
	return c.userPresenter.PresentUpdateUserName(e, output)
}
