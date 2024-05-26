package createuser_test

import (
	"fmt"
	"testing"

	status "github.com/Go_CleanArch/common/const"
	createUserDomain "github.com/Go_CleanArch/domain/factory/user/create_user"
	"github.com/stretchr/testify/assert"
)

// generatePasswordのモック関数
func mockGeneratePassword(password string) (string, error) {
	return "", fmt.Errorf("encryption failed %s", password)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	email := "test@example.com"
	userName := "Test User"
	password := "password123"

	t.Run("正常系: 有効なプロパティでユーザー作成", func(t *testing.T) {
		t.Parallel()

		factory := createUserDomain.NewCreateUserFactory()
		user, err := factory.CreateUser(&createUserDomain.CreateUserInitProps{
			UserId:   "",
			UserName: userName,
			Password: password,
			Email:    email,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, user.UserID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, userName, user.UserName)
		assert.NotEqual(t, password, user.Password)
	})

	t.Run("異常系: ユーザーIDが既に存在する", func(t *testing.T) {
		t.Parallel()

		existingUserId := "existing_user_id"
		factory := createUserDomain.NewCreateUserFactory()
		user, err := factory.CreateUser(&createUserDomain.CreateUserInitProps{
			UserId:   existingUserId,
			UserName: userName,
			Password: password,
			Email:    email,
		})

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusCode, err.Status)
		assert.Contains(t, err.Messages[0].Key, "email")
		assert.Contains(t, err.Messages[0].Value, "すでに登録されているアドレスです")
		assert.Contains(t, err.Detail, "Enable Check Error")
	})

	t.Run("異常系: パスワードの暗号化に失敗", func(t *testing.T) {
		t.Parallel()

		// CreateUserFactoryのフィールドにモック関数を注入
		factory := &createUserDomain.CreateUserFactory{
			GeneratePassword: mockGeneratePassword,
		}

		user, err := factory.CreateUser(&createUserDomain.CreateUserInitProps{
			UserId:   "",
			UserName: userName,
			Password: password,
			Email:    email,
		})

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusCode, err.Status)
		assert.Contains(t, err.Detail, "Internal Server Error")
	})
}