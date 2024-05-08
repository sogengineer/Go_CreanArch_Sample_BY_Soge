package user_service_impl_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Go_CleanArch/common/crypto"

	"github.com/Go_CleanArch/interface_adapter/gateway/entity"
	inputUser "github.com/Go_CleanArch/usecase/input/user"
	user_service_impl "github.com/Go_CleanArch/usecase/service/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, userJson []byte) (*entity.User, error) {
	args := m.Called(ctx, userJson)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUserService(t *testing.T) {
	t.Parallel()
	t.Run("新規ユーザー作成_正常系", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mockUserRepo := new(MockUserRepository)
		userService := user_service_impl.NewUserService(mockUserRepo)
		// テスト用のリクエストボディを作成
		createUserForm := inputUser.CreateUserForm{
			Email:    "test@example.com",
			UserName: "testuser",
			Password: "Password123",
		}
		requestBody, _ := json.Marshal(createUserForm)

		// モックの設定
		mockUserRepo.On("FindUserByEmail", ctx, createUserForm.Email).Return(&entity.User{}, nil)
		mockUserRepo.On("CreateUser", ctx, mock.Anything).Return(&entity.User{UserId: "user123"}, nil)

		// リクエストの作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		presenter, err := userService.CreateUserService(ctx, c)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, "user123", presenter.UserId)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("新規ユーザー作成_バリデーションエラー", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mockUserRepo := new(MockUserRepository)
		userService := user_service_impl.NewUserService(mockUserRepo)
		// テスト用のリクエストボディを作成
		createUserForm := inputUser.CreateUserForm{
			Email:    "invalid_email",
			UserName: "",
			Password: "password",
		}
		requestBody, _ := json.Marshal(createUserForm)

		// リクエストの作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		_, err := userService.CreateUserService(ctx, c)

		// アサーション
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "Validation error occurred")
	})

	t.Run("新規ユーザー作成_メールアドレス重複", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mockUserRepo := new(MockUserRepository)
		userService := user_service_impl.NewUserService(mockUserRepo)
		// テスト用のリクエストボディを作成
		createUserForm := inputUser.CreateUserForm{
			Email:    "test@example.com",
			UserName: "testuser",
			Password: "Password123",
		}
		requestBody, _ := json.Marshal(createUserForm)

		// モックの設定
		mockUserRepo.On("FindUserByEmail", ctx, createUserForm.Email).Return(&entity.User{UserId: "user123"}, nil)

		// リクエストの作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		_, err := userService.CreateUserService(ctx, c)

		// アサーション
		assert.Error(t, err)
		// assert.Contains(t, err.Error(), "Failed to build user factory props")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestLoginService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	userService := user_service_impl.NewUserService(mockUserRepo)

	t.Run("ログイン_正常系", func(t *testing.T) {
		t.Parallel()
		// テスト用のリクエストボディを作成
		loginForm := inputUser.LoginForm{
			Email:    "test@example.com",
			Password: "password",
		}
		requestBody, _ := json.Marshal(loginForm)

		// モックの設定
		hashedPassword, _ := crypto.PasswordEncrypt("password")
		mockUserRepo.On("FindUserByEmail", ctx, loginForm.Email).Return(&entity.User{
			Email:    loginForm.Email,
			UserName: "testuser",
			Password: hashedPassword,
		}, nil)

		// リクエストの作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		presenter, err := userService.LoginService(ctx, c)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, loginForm.Email, presenter.Email)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ログイン_バリデーションエラー", func(t *testing.T) {
		t.Parallel()
		// テスト用のリクエストボディを作成
		loginForm := inputUser.LoginForm{
			Email:    "invalid_email",
			Password: "password",
		}
		requestBody, _ := json.Marshal(loginForm)

		// リクエストの作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		result, err := userService.LoginService(ctx, c)

		// アサーション
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUserRepo.AssertNotCalled(t, "FindUserByEmail")
		mockUserRepo.AssertNotCalled(t, "CreateUser")
	})

	t.Run("ログイン_メールアドレス存在しない", func(t *testing.T) {
		t.Parallel()
		// テスト用のリクエストボディを作成
		loginForm := inputUser.LoginForm{
			Email:    "nonexistent@example.com",
			Password: "Password1234",
		}
		requestBody, _ := json.Marshal(loginForm)

		// モックの設定
		mockUserRepo.On("FindUserByEmail", ctx, loginForm.Email).Return(&entity.User{}, fmt.Errorf("条件に一致するレコードが見つかりません: record not found"))

		// リクエストの作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))

		// テスト対象の関数を実行
		_, err := userService.LoginService(ctx, c)
		// アサーション
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Emailアドレスが存在しません")
		mockUserRepo.AssertExpectations(t)
	})
}
