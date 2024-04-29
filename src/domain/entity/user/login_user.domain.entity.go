package login_user_domain_entity

import (
	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/common/errors"
)

// LoginUserBuilderProps はイベント作成のプロパティを定義する。
type LoginUserBuilderProps struct {
	UserId   string
	UserName string
	Email    string
}

type LoginUserBuilderPropsOption func(*LoginUserBuilderProps) ([]errors.ApiErrMessage, error)

func NewLoginUserDomainProps(opts ...LoginUserBuilderPropsOption) (*LoginUserBuilderProps, *errors.ApiErr) {
	apiErrMessages := make([]errors.ApiErrMessage, 0)
	presenter := &LoginUserBuilderProps{}
	for _, opt := range opts {
		setErrMessages, err := opt(presenter)
		if err != nil {
			// エラーが発生した場合は、Internal Server Error を返す
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

	return presenter, nil
}

func WithLoginUserIdAndEmail(email string) LoginUserBuilderPropsOption {
	return func(props *LoginUserBuilderProps) ([]errors.ApiErrMessage, error) {
		props.UserId = crypto.GenerateUserId(email)
		props.Email = email
		return nil, nil
	}
}

func WithLoginUserName(userName string) LoginUserBuilderPropsOption {
	return func(props *LoginUserBuilderProps) ([]errors.ApiErrMessage, error) {
		props.UserName = userName
		return nil, nil
	}
}

func WithLoginPassword(targetPassword string, sourcePassword string) LoginUserBuilderPropsOption {
	return func(props *LoginUserBuilderProps) ([]errors.ApiErrMessage, error) {
		err := crypto.CompareHashAndPassword(targetPassword, sourcePassword)
		if err != nil {
			return []errors.ApiErrMessage{
				{
					Key:   "password",
					Value: "パスワードが間違っています",
				},
			}, nil
		}
		return nil, nil
	}
}
