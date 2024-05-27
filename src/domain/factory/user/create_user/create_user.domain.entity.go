package createuser

import (
	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/common/errors"
	entity "github.com/Go_CleanArch/domain/entity"
	log "github.com/sirupsen/logrus"
)

type CreateUserInitProps struct {
	UserId   string
	UserName string
	Password string
	Email    string
}

type CreateUserFactory struct {
	GeneratePassword func(string) (string, error)
}

func NewCreateUserFactory() *CreateUserFactory {
	return &CreateUserFactory{
		GeneratePassword: generatePassword,
	}
}

func (uf *CreateUserFactory) CreateUser(props *CreateUserInitProps) (*entity.User, *errors.ApiErr) {
	apiErrMessages := make([]errors.ApiErrMessage, 0)
	// ユーザーIDの存在チェック
	checkUserExistErrorMessage, err := checkUserExist(props.UserId)
	if err != nil {
		log.WithError(err).Error("Failed to check user existence")
		return nil, errors.OutputApiError(
			[]errors.ApiErrMessage{
				{
					Key:   "undefined",
					Value: err.Error(),
				},
			},
			status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusCode,
			status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusName,
		)
	}
	// checkUserExistErrorMessageの値がnil以外だった場合APIエラーメッセージリストに内容を追加する
	if checkUserExistErrorMessage != nil {
		apiErrMessages = append(apiErrMessages, checkUserExistErrorMessage...)
	}

	// パスワードのハッシュ化
	hashedPassword, err := uf.GeneratePassword(props.Password)
	if err != nil {
		log.WithError(err).Error("Failed to generate password")
		return nil, errors.OutputApiError(
			[]errors.ApiErrMessage{
				{
					Key:   "undefined",
					Value: err.Error(),
				},
			},
			status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusCode,
			status.ErrorStatusMap["INTERNAL_SERVER_ERROR"].StatusName,
		)
	}

	// entity.Userを生成して返す
	user, newUserErrorMessage := entity.NewUser(
		entity.WithUserID(crypto.GenerateUserId(props.Email)),
		entity.WithUserName(props.UserName),
		entity.WithPassword(hashedPassword),
		entity.WithEmail(props.Email),
	)
	if newUserErrorMessage != nil {
		apiErrMessages = append(apiErrMessages, newUserErrorMessage.Messages...)
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

func checkUserExist(userId string) ([]errors.ApiErrMessage, error) {
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

func generatePassword(password string) (string, error) {
	hashedPw, err := crypto.PasswordEncrypt(password)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt password")
		return "", err
	}
	return hashedPw, nil
}
