package create_user_domain_entity

import (
	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/common/errors"
	log "github.com/sirupsen/logrus"
)

// CreateUserFactoryProps は、ユーザー作成に必要なプロパティを持つ構造体
type CreateUserFactoryProps struct {
	UserId   string
	UserName string
	Password string
	Email    string
}

// CreateUserFactoryPropsOption は、CreateUserFactoryProps を変更するための関数オプション
type CreateUserFactoryPropsOption func(*CreateUserFactoryProps) ([]errors.ApiErrMessage, error)

// NewCreateUserFactoryProps は、与えられたオプションを適用して新しい CreateUserFactoryProps を作成する
func NewCreateUserFactoryProps(opts ...CreateUserFactoryPropsOption) (*CreateUserFactoryProps, *errors.ApiErr) {
	apiErrMessages := make([]errors.ApiErrMessage, 0)
	createUserFactoryProps := &CreateUserFactoryProps{}

	// 各オプションを適用
	for _, opt := range opts {
		setErrMessages, err := opt(createUserFactoryProps)
		if err != nil {
			// エラーが発生した場合は、Internal Server Error を返す
			log.WithError(err).Error("Failed to apply user factory props option")
			apiErr := errors.OutputApiError(
				[]errors.ApiErrMessage{
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
		log.WithField("apiErrMessages", apiErrMessages).Error("Validation errors occurred")
		apiErr := errors.OutputApiError(
			apiErrMessages,
			status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusCode,
			status.ErrorStatusMap["ENABLE_CHECK_ERROR"].StatusName,
		)
		return nil, apiErr
	}

	log.WithField("userId", createUserFactoryProps.UserId).Info("User factory props created successfully")
	return createUserFactoryProps, nil
}

// WithUserId は、ユーザーIDの存在チェックを行い、存在していた場合有効性チェックエラーメッセージをreturnするオプション
func WithUserId(userId string) CreateUserFactoryPropsOption {
	return func(props *CreateUserFactoryProps) ([]errors.ApiErrMessage, error) {
		if userId != "" {
			// ユーザーIDが空でない場合は、エラーメッセージを返す
			log.WithField("userId", userId).Info("User ID already exists")
			return []errors.ApiErrMessage{
				{
					Key:   "email",
					Value: "すでに登録されているアドレスです",
				},
			}, nil
		}
		return nil, nil
	}
}

// WithEmail は、メールアドレスからユーザーIDを生成し、プロパティにセットするオプション
func WithEmail(email string) CreateUserFactoryPropsOption {
	return func(props *CreateUserFactoryProps) ([]errors.ApiErrMessage, error) {
		props.UserId = crypto.GenerateUserId(email)
		props.Email = email
		return nil, nil
	}
}

// WithUserName は、ユーザー名をプロパティにセットするオプション
func WithUserName(userName string) CreateUserFactoryPropsOption {
	return func(props *CreateUserFactoryProps) ([]errors.ApiErrMessage, error) {
		props.UserName = userName
		log.WithField("userName", userName).Info("User name set")
		return nil, nil
	}
}

// WithPassword は、パスワードをハッシュ化してプロパティにセットするオプション
func WithPassword(password string) CreateUserFactoryPropsOption {
	return func(props *CreateUserFactoryProps) ([]errors.ApiErrMessage, error) {
		hashedPw, err := crypto.PasswordEncrypt(password)
		if err != nil {
			log.WithError(err).Error("Failed to encrypt password")
			return nil, err
		}
		props.Password = hashedPw
		log.Info("Password encrypted and set")
		return nil, nil
	}
}
