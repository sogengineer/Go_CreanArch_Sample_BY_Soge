package create_user_domain_service_test

import (
	"fmt"
	"testing"

	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/errors"
	createUserDomain "github.com/Go_CleanArch/domain/service/user/create_user"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateUserDomainServiceProps(t *testing.T) {
	t.Parallel()

	email := "test@example.com"
	userName := "Test User"
	password := "password123"

	t.Run("正常系: 有効なプロパティでユーザー作成", func(t *testing.T) {
		t.Parallel()

		props, err := createUserDomain.NewCreateUserDomainServiceProps(
			createUserDomain.WithUserId(""),
			createUserDomain.WithEmail(email),
			createUserDomain.WithUserName(userName),
			createUserDomain.WithPassword(password),
		)

		assert.Nil(t, err)
		assert.NotEmpty(t, props.UserId)
		assert.Equal(t, email, props.Email)
		assert.Equal(t, userName, props.UserName)
		assert.NotEqual(t, password, props.Password)
	})

	t.Run("異常系: ユーザーIDが既に存在する", func(t *testing.T) {
		t.Parallel()

		existingUserId := "existing_user_id"

		props, err := createUserDomain.NewCreateUserDomainServiceProps(
			createUserDomain.WithUserId(existingUserId),
			createUserDomain.WithEmail(email),
			createUserDomain.WithUserName(userName),
			createUserDomain.WithPassword(password),
		)

		assert.Nil(t, props)
		assert.NotNil(t, err)
		assert.Equal(t, status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusCode, err.Status)
		assert.Contains(t, err.Messages[0].Key, "email")
		assert.Contains(t, err.Messages[0].Value, "すでに登録されているアドレスです")
		assert.Contains(t, err.Detail, "Enable Check Error")
	})

	t.Run("異常系: パスワードの暗号化に失敗", func(t *testing.T) {
		t.Parallel()

		withPasswordMock := func(password string) createUserDomain.CreateUserDomainServicePropsOption {
			return func(props *createUserDomain.CreateUserDomainServiceProps) ([]errors.ApiErrMessage, error) {
				return nil, fmt.Errorf("encryption failed %s", password)
			}
		}

		props, err := createUserDomain.NewCreateUserDomainServiceProps(
			createUserDomain.WithUserId(""),
			createUserDomain.WithEmail(email),
			createUserDomain.WithUserName(userName),
			withPasswordMock(""),
		)

		assert.Nil(t, props)
		assert.NotNil(t, err)
		assert.Equal(t, status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusCode, err.Status)
		assert.Contains(t, err.Detail, "Internal Server Error")
	})
}
