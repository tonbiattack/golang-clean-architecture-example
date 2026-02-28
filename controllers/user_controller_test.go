// Package controllers_test はcontrollersパッケージのテストを提供します。
package controllers_test

import (
	"context"
	"errors"
	"golang-clean-architecture-example/controllers"
	"golang-clean-architecture-example/domain/entities"
	"golang-clean-architecture-example/usecases/dto/input"
	"golang-clean-architecture-example/usecases/dto/output"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// mockUpdateUserNameInteractor はテスト用のモックインタラクターです。
type mockUpdateUserNameInteractor struct {
	executeFunc func(ctx context.Context, input *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error)
}

func (m *mockUpdateUserNameInteractor) Execute(ctx context.Context, in *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, in)
	}
	return nil, errors.New("Execute not implemented")
}

// mockUserPresenter はテスト用のモックプレゼンターです。
type mockUserPresenter struct {
	presentUpdateUserNameFunc func(c echo.Context, output *output.UpdateUserNameOutput) error
}

func (m *mockUserPresenter) PresentUpdateUserName(c echo.Context, out *output.UpdateUserNameOutput) error {
	if m.presentUpdateUserNameFunc != nil {
		return m.presentUpdateUserNameFunc(c, out)
	}
	return nil
}

// mockErrorPresenter はテスト用のモックエラープレゼンターです。
type mockErrorPresenter struct {
	presentBadRequestFunc          func(c echo.Context, message string) error
	presentInternalServerErrorFunc func(c echo.Context, err error) error
}

func (m *mockErrorPresenter) PresentBadRequest(c echo.Context, message string) error {
	if m.presentBadRequestFunc != nil {
		return m.presentBadRequestFunc(c, message)
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
}

func (m *mockErrorPresenter) PresentInternalServerError(c echo.Context, err error) error {
	if m.presentInternalServerErrorFunc != nil {
		return m.presentInternalServerErrorFunc(c, err)
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

// TestNewUserController はNewUserControllerが正しくインスタンスを生成することをテストします。
func TestNewUserController(t *testing.T) {
	mockInteractor := &mockUpdateUserNameInteractor{}
	mockUserPres := &mockUserPresenter{}
	mockErrPres := &mockErrorPresenter{}

	controller := controllers.NewUserController(mockInteractor, mockUserPres, mockErrPres)

	if controller == nil {
		t.Fatal("NewUserController returned nil")
	}
}

// TestUserController_UpdateUserName_Success は正常系のテストです。
func TestUserController_UpdateUserName_Success(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		reqBody  string
		wantName string
	}{
		{
			name:     "通常のユーザー名更新",
			userID:   "user001",
			reqBody:  `{"name": "New Name"}`,
			wantName: "New Name",
		},
		{
			name:     "日本語のユーザー名更新",
			userID:   "user002",
			reqBody:  `{"name": "新しい名前"}`,
			wantName: "新しい名前",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			mockInteractor := &mockUpdateUserNameInteractor{
				executeFunc: func(ctx context.Context, in *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error) {
					if in.UserID != tt.userID {
						t.Errorf("Execute called with wrong userID: got %v, want %v", in.UserID, tt.userID)
					}
					if in.NewName != tt.wantName {
						t.Errorf("Execute called with wrong newName: got %v, want %v", in.NewName, tt.wantName)
					}
					return &output.UpdateUserNameOutput{
						User: entities.NewUser(tt.userID, tt.wantName),
					}, nil
				},
			}

			presentCalled := false
			mockUserPres := &mockUserPresenter{
				presentUpdateUserNameFunc: func(c echo.Context, out *output.UpdateUserNameOutput) error {
					presentCalled = true
					if out.User.GetID() != tt.userID {
						t.Errorf("PresentUpdateUserName called with wrong user ID: got %v, want %v", out.User.GetID(), tt.userID)
					}
					if out.User.GetName() != tt.wantName {
						t.Errorf("PresentUpdateUserName called with wrong user name: got %v, want %v", out.User.GetName(), tt.wantName)
					}
					return c.JSON(http.StatusOK, map[string]string{
						"id":   out.User.GetID(),
						"name": out.User.GetName(),
					})
				},
			}

			mockErrPres := &mockErrorPresenter{}

			controller := controllers.NewUserController(mockInteractor, mockUserPres, mockErrPres)

			// HTTPリクエストのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID+"/update_name", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			err := controller.UpdateUserName(c, tt.userID)

			// 検証
			if err != nil {
				t.Errorf("UpdateUserName() returned unexpected error: %v", err)
			}

			if !presentCalled {
				t.Error("PresentUpdateUserName was not called")
			}

			if rec.Code != http.StatusOK {
				t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

// TestUserController_UpdateUserName_InvalidJSON は無効なJSONの場合のテストです。
func TestUserController_UpdateUserName_InvalidJSON(t *testing.T) {
	mockInteractor := &mockUpdateUserNameInteractor{}
	mockUserPres := &mockUserPresenter{}

	badRequestCalled := false
	mockErrPres := &mockErrorPresenter{
		presentBadRequestFunc: func(c echo.Context, message string) error {
			badRequestCalled = true
			if message != "invalid request" {
				t.Errorf("PresentBadRequest called with wrong message: got %v, want 'invalid request'", message)
			}
			return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
		},
	}

	controller := controllers.NewUserController(mockInteractor, mockUserPres, mockErrPres)

	// 無効なJSONでリクエストを作成
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/user001/update_name", strings.NewReader(`{invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	err := controller.UpdateUserName(c, "user001")

	// 検証
	if err != nil {
		t.Errorf("UpdateUserName() returned unexpected error: %v", err)
	}

	if !badRequestCalled {
		t.Error("PresentBadRequest was not called")
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestUserController_UpdateUserName_InteractorError はインタラクターがエラーを返す場合のテストです。
func TestUserController_UpdateUserName_InteractorError(t *testing.T) {
	expectedError := errors.New("database error")

	mockInteractor := &mockUpdateUserNameInteractor{
		executeFunc: func(ctx context.Context, in *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error) {
			return nil, expectedError
		},
	}

	mockUserPres := &mockUserPresenter{}

	internalServerErrorCalled := false
	mockErrPres := &mockErrorPresenter{
		presentInternalServerErrorFunc: func(c echo.Context, err error) error {
			internalServerErrorCalled = true
			if err != expectedError {
				t.Errorf("PresentInternalServerError called with wrong error: got %v, want %v", err, expectedError)
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		},
	}

	controller := controllers.NewUserController(mockInteractor, mockUserPres, mockErrPres)

	// HTTPリクエストのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/user001/update_name", strings.NewReader(`{"name": "New Name"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	err := controller.UpdateUserName(c, "user001")

	// 検証
	if err != nil {
		t.Errorf("UpdateUserName() returned unexpected error: %v", err)
	}

	if !internalServerErrorCalled {
		t.Error("PresentInternalServerError was not called")
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// TestUserController_UpdateUserName_EmptyName は空の名前での更新をテストします。
func TestUserController_UpdateUserName_EmptyName(t *testing.T) {
	mockInteractor := &mockUpdateUserNameInteractor{
		executeFunc: func(ctx context.Context, in *input.UpdateUserNameInput) (*output.UpdateUserNameOutput, error) {
			if in.NewName != "" {
				t.Errorf("Execute called with non-empty name: got %v", in.NewName)
			}
			return &output.UpdateUserNameOutput{
				User: entities.NewUser("user001", ""),
			}, nil
		},
	}

	mockUserPres := &mockUserPresenter{
		presentUpdateUserNameFunc: func(c echo.Context, out *output.UpdateUserNameOutput) error {
			return c.JSON(http.StatusOK, map[string]string{
				"id":   out.User.GetID(),
				"name": out.User.GetName(),
			})
		},
	}

	mockErrPres := &mockErrorPresenter{}

	controller := controllers.NewUserController(mockInteractor, mockUserPres, mockErrPres)

	// HTTPリクエストのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/user001/update_name", strings.NewReader(`{"name": ""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	err := controller.UpdateUserName(c, "user001")

	// 検証
	if err != nil {
		t.Errorf("UpdateUserName() returned unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}
