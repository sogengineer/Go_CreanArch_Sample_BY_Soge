// domain/entity/user.go
package entity

import (
	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/errors"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserID   string
	UserName string
	Password string
	Email    string
}

type UserOption func(*User) ([]errors.ApiErrMessage, error)

func NewUser(opts ...UserOption) (*User, *errors.ApiErr) {
	apiErrMessages := make([]errors.ApiErrMessage, 0)
	user := &User{}
	for _, opt := range opts {
		setErrMessages, err := opt(user)
		if err != nil {
			// エラーが発生した場合は、Internal Server Error を返す
			log.WithError(err).Error("login user domain entity INTERNAL_SERVER_ERROR")
			apiErr := errors.OutputApiError([]errors.ApiErrMessage{
				{
					Key:   "undefined",
					Value: err.Error(),
				},
			},
				status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusCode,
				status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusName,
			)
			return nil, apiErr
		}
		apiErrMessages = append(apiErrMessages, setErrMessages...)
	}

	// エラーメッセージがある場合は、EnableCheckError を返す
	if len(apiErrMessages) > 0 {
		apiErr := errors.OutputApiError(
			apiErrMessages,
			status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusCode,
			status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusName,
		)
		return nil, apiErr
	}

	return user, nil
}

func WithUserID(userID string) UserOption {
	return func(u *User) ([]errors.ApiErrMessage, error) {
		u.UserID = userID
		return nil, nil
	}
}

func WithUserName(userName string) UserOption {
	return func(u *User) ([]errors.ApiErrMessage, error) {
		u.UserName = userName
		return nil, nil
	}
}

func WithPassword(password string) UserOption {
	return func(u *User) ([]errors.ApiErrMessage, error) {
		u.Password = password
		return nil, nil
	}
}

func WithEmail(email string) UserOption {
	return func(u *User)  ([]errors.ApiErrMessage, error) {
		u.Email = email
		return nil, nil
	}
}
