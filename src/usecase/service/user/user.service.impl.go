package user_service_impl

import (
	"context"

	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/common/errors"
	loginUserDomain "github.com/Go_CleanArch/domain/entity/user"
	createUserDomain "github.com/Go_CleanArch/domain/factory/user"
	inputUser "github.com/Go_CleanArch/usecase/input/user"
	outputUser "github.com/Go_CleanArch/usecase/output/user"
	repository "github.com/Go_CleanArch/usecase/repository_interface"
	"github.com/gin-gonic/gin"
)

// Service provides user's behavior
type UserService struct {
	userRepository repository.UserRepository
}

// Constructor
func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// サインアップ
func (us *UserService) CreateUserService(ctx context.Context, c *gin.Context) (outputUser.CreateUserPresenter, error) {
	var createUserForm inputUser.CreateUserForm
	var createUserPresenter outputUser.CreateUserPresenter
	if err := c.BindJSON(&createUserForm); err != nil {
		return createUserPresenter, err
	}

	// 入力チェックバリデーション
	apiErrMessages := createUserForm.CreateUserValidate()
	if len(apiErrMessages) > 0 {
		apiErr := errors.OutputApiError(
			apiErrMessages,
			status.ErrorStatusMap["BAD_REQUEST"].StatusCode,
			status.ErrorStatusMap["BAD_REQUEST"].StatusName,
		)
		c.JSON(apiErr.Status, apiErr)
		return createUserPresenter, apiErr.Error()
	}

	// 登録済みのメールアドレスを再登録しようとしていないかチェック
	findUser, _ := us.userRepository.FindUserByEmail(ctx, createUserForm.Email)
	findUserId := ""
	if findUser != nil {
		findUserId = findUser.UserId
	}

	// 登録するユーザー情報のビルドを行う
	createUserFactoryProps, apiErr := createUserDomain.NewCreateUserFactoryProps(
		createUserDomain.WithUserId(findUserId),
		createUserDomain.WithEmail(createUserForm.Email),
		createUserDomain.WithUserName(createUserForm.UserName),
		createUserDomain.WithPassword(createUserForm.Password),
	)
	if apiErr != nil {
		c.JSON(apiErr.Status, apiErr)
		return createUserPresenter, apiErr.Error()
	}

	// ビルドしたユーザー情報を基にユーザー登録を行う
	getUserJson, err := crypto.ConvertStructIntoJson(createUserFactoryProps)
	if err != nil {
		c.JSON(500, err)
		return createUserPresenter, err
	}
	createdUser, err := us.userRepository.CreateUser(ctx, getUserJson)
	if err != nil {
		c.JSON(500, err)
		return createUserPresenter, err
	}
	if err := crypto.ConvertJsonAndCopyBean(createdUser, &createUserPresenter); err != nil {
		return createUserPresenter, err
	}
	return createUserPresenter, nil
}

// ログイン
func (us *UserService) LoginService(ctx context.Context, c *gin.Context) (outputUser.LoginPresenter, error) {
	var loginForm inputUser.LoginForm
	var loginPresenter outputUser.LoginPresenter

	if err := c.BindJSON(&loginForm); err != nil {
		return loginPresenter, err
	}

	// 入力チェックバリデーション
	apiErrMessages := loginForm.LoginValidate()
	if len(apiErrMessages) > 0 {
		apiErr := errors.OutputApiError(
			apiErrMessages,
			status.ErrorStatusMap["BAD_REQUEST"].StatusCode,
			status.ErrorStatusMap["BAD_REQUEST"].StatusName,
		)
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	getUser, err := us.userRepository.FindUserByEmail(ctx, loginForm.Email)
	// メールアドレス確認
	if err != nil {
		apiErr := errors.OutputApiError(append(apiErrMessages,
			errors.ApiErrMessage{
				Key:   "email",            // エラーの発生したフィールド名
				Value: "Emailアドレスが存在しません", // エラーメッセージ
			}),
			status.ErrorStatusMap["NOT_FOUND"].StatusCode,
			status.ErrorStatusMap["NOT_FOUND"].StatusName,
		)
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	loginUserDomainEntity, apiErr := loginUserDomain.NewLoginUserDomainProps(
		loginUserDomain.WithLoginUserIdAndEmail(getUser.Email),
		loginUserDomain.WithLoginUserName(getUser.UserName),
		loginUserDomain.WithLoginPassword(getUser.Password, loginForm.Password),
	)
	if apiErr != nil {
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	getUserJson, err := crypto.ConvertStructIntoJson(loginUserDomainEntity)
	if err != nil {
		return loginPresenter, err
	}
	if err := crypto.CopyBeans(getUserJson, &loginPresenter); err != nil {
		return loginPresenter, err
	}

	return loginPresenter, nil
}
