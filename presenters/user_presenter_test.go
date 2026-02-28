// Package presenters_test はpresentersパッケージのテストを提供します。
package presenters_test

import (
	"golang-clean-architecture-example/domain/entities"
	"golang-clean-architecture-example/presenters"
	"golang-clean-architecture-example/usecases/dto/output"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestNewUserPresenter はNewUserPresenterが正しくインスタンスを生成することをテストします。
func TestNewUserPresenter(t *testing.T) {
	presenter := presenters.NewUserPresenter()

	if presenter == nil {
		t.Fatal("NewUserPresenter returned nil")
	}
}

// TestUserPresenter_PresentUpdateUserName はPresentUpdateUserNameの動作をテストします。
func TestUserPresenter_PresentUpdateUserName(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		userName string
	}{
		{
			name:     "通常のレスポンス",
			userID:   "user001",
			userName: "Test User",
		},
		{
			name:     "日本語を含むレスポンス",
			userID:   "user002",
			userName: "テストユーザー",
		},
		{
			name:     "空文字を含むレスポンス",
			userID:   "user003",
			userName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プレゼンターを作成
			presenter := presenters.NewUserPresenter()

			// 出力DTOを作成
			outputDTO := &output.UpdateUserNameOutput{
				User: entities.NewUser(tt.userID, tt.userName),
			}

			// Echoコンテキストのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// PresentUpdateUserNameを実行
			err := presenter.PresentUpdateUserName(c, outputDTO)

			// 検証
			if err != nil {
				t.Errorf("PresentUpdateUserName() returned unexpected error: %v", err)
			}

			if rec.Code != http.StatusOK {
				t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
			}

			// Content-Typeの確認（application/jsonまたはapplication/json; charset=UTF-8）
			contentType := rec.Header().Get(echo.HeaderContentType)
			if contentType != echo.MIMEApplicationJSON && contentType != echo.MIMEApplicationJSONCharsetUTF8 {
				t.Errorf("Content-Type = %v, want %v or %v", contentType, echo.MIMEApplicationJSON, echo.MIMEApplicationJSONCharsetUTF8)
			}

			// レスポンスボディが空でないことを確認
			if rec.Body.Len() == 0 {
				t.Error("Response body is empty")
			}
		})
	}
}

// TestNewErrorPresenter はNewErrorPresenterが正しくインスタンスを生成することをテストします。
func TestNewErrorPresenter(t *testing.T) {
	presenter := presenters.NewErrorPresenter()

	if presenter == nil {
		t.Fatal("NewErrorPresenter returned nil")
	}
}

// TestErrorPresenter_PresentBadRequest はPresentBadRequestの動作をテストします。
func TestErrorPresenter_PresentBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "通常のエラーメッセージ",
			message: "invalid request",
		},
		{
			name:    "日本語のエラーメッセージ",
			message: "無効なリクエストです",
		},
		{
			name:    "空のエラーメッセージ",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プレゼンターを作成
			presenter := presenters.NewErrorPresenter()

			// Echoコンテキストのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// PresentBadRequestを実行
			err := presenter.PresentBadRequest(c, tt.message)

			// 検証
			if err != nil {
				t.Errorf("PresentBadRequest() returned unexpected error: %v", err)
			}

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Status code = %d, want %d", rec.Code, http.StatusBadRequest)
			}

			// Content-Typeの確認（application/jsonまたはapplication/json; charset=UTF-8）
			contentType := rec.Header().Get(echo.HeaderContentType)
			if contentType != echo.MIMEApplicationJSON && contentType != echo.MIMEApplicationJSONCharsetUTF8 {
				t.Errorf("Content-Type = %v, want %v or %v", contentType, echo.MIMEApplicationJSON, echo.MIMEApplicationJSONCharsetUTF8)
			}

			// レスポンスボディが空でないことを確認
			if rec.Body.Len() == 0 {
				t.Error("Response body is empty")
			}
		})
	}
}

// TestErrorPresenter_PresentInternalServerError はPresentInternalServerErrorの動作をテストします。
func TestErrorPresenter_PresentInternalServerError(t *testing.T) {
	tests := []struct {
		name      string
		errorMsg  string
	}{
		{
			name:      "通常のエラー",
			errorMsg:  "database connection error",
		},
		{
			name:      "日本語のエラー",
			errorMsg:  "データベース接続エラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プレゼンターを作成
			presenter := presenters.NewErrorPresenter()

			// エラーを作成
			testError := &testError{message: tt.errorMsg}

			// Echoコンテキストのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// PresentInternalServerErrorを実行
			err := presenter.PresentInternalServerError(c, testError)

			// 検証
			if err != nil {
				t.Errorf("PresentInternalServerError() returned unexpected error: %v", err)
			}

			if rec.Code != http.StatusInternalServerError {
				t.Errorf("Status code = %d, want %d", rec.Code, http.StatusInternalServerError)
			}

			// Content-Typeの確認（application/jsonまたはapplication/json; charset=UTF-8）
			contentType := rec.Header().Get(echo.HeaderContentType)
			if contentType != echo.MIMEApplicationJSON && contentType != echo.MIMEApplicationJSONCharsetUTF8 {
				t.Errorf("Content-Type = %v, want %v or %v", contentType, echo.MIMEApplicationJSON, echo.MIMEApplicationJSONCharsetUTF8)
			}

			// レスポンスボディが空でないことを確認
			if rec.Body.Len() == 0 {
				t.Error("Response body is empty")
			}
		})
	}
}

// testError はテスト用のエラー構造体です。
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
