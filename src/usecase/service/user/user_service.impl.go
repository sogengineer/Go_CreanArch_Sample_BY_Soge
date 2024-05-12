package user_service_impl

import (
	"context"

	status "github.com/Go_CleanArch/common/const"
	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/common/errors"
	createUserDomainService "github.com/Go_CleanArch/domain/service/user/create_user"
	loginUserDomainService "github.com/Go_CleanArch/domain/service/user/login_user"
	inputUser "github.com/Go_CleanArch/usecase/input/user"
	outputUser "github.com/Go_CleanArch/usecase/output/user"
	repository "github.com/Go_CleanArch/usecase/repository_interface"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Error("Failed to bind JSON request body")
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
		log.WithField("apiErr", apiErr).Error("Validation error occurred")
		c.JSON(apiErr.Status, apiErr)
		return createUserPresenter, apiErr.Error()
	}

	// 登録済みのメールアドレスを再登録しようとしていないかチェック
	findUser, err := us.userRepository.FindUserByEmail(ctx, createUserForm.Email)
	if err != nil {
		log.WithError(err).Error("Failed to find user by email")
	}
	findUserId := ""
	if findUser != nil {
		findUserId = findUser.UserId
	}

	// 登録するユーザー情報のビルドを行う
	createUserDomainServiceProps, apiErr := createUserDomainService.NewCreateUserDomainServiceProps(
		createUserDomainService.WithUserId(findUserId),
		createUserDomainService.WithEmail(createUserForm.Email),
		createUserDomainService.WithUserName(createUserForm.UserName),
		createUserDomainService.WithPassword(createUserForm.Password),
	)
	if apiErr != nil {
		log.WithField("apiErr", apiErr).Error("Failed to build user factory props")
		c.JSON(apiErr.Status, apiErr)
		return createUserPresenter, apiErr.Error()
	}

	// ビルドしたユーザー情報を基にユーザー登録を行う
	getUserJson, err := crypto.ConvertStructIntoJson(createUserDomainServiceProps)
	if err != nil {
		log.WithError(err).Error("Failed to convert user factory props into JSON")
		c.JSON(500, err)
		return createUserPresenter, err
	}
	createdUser, err := us.userRepository.CreateUser(ctx, getUserJson)
	if err != nil {
		log.WithError(err).Error("Failed to create user")
		c.JSON(500, err)
		return createUserPresenter, err
	}
	if err := crypto.ConvertJsonAndCopyBean(createdUser, &createUserPresenter); err != nil {
		log.WithError(err).Error("Failed to convert created user into presenter")
		return createUserPresenter, err
	}
	log.WithField("userId", createUserPresenter.UserId).Info("User created successfully")
	return createUserPresenter, nil
}

// ログイン
func (us *UserService) LoginService(ctx context.Context, c *gin.Context) (outputUser.LoginPresenter, error) {
	var loginForm inputUser.LoginForm
	var loginPresenter outputUser.LoginPresenter

	if err := c.BindJSON(&loginForm); err != nil {
		log.WithError(err).Error("Failed to bind JSON request body")
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
		log.WithField("apiErr", apiErr).Error("Validation error occurred")
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	getUser, err := us.userRepository.FindUserByEmail(ctx, loginForm.Email)
	// メールアドレス確認
	if err != nil {
		log.WithError(err).Error("Failed to find user by email")
		apiErr := errors.OutputApiError(
			append(apiErrMessages,
				errors.ApiErrMessage{
					Key:   "email",
					Value: "Emailアドレスが存在しません",
				},
			),
			status.ErrorStatusMap["NOT_FOUND"].StatusCode,
			status.ErrorStatusMap["NOT_FOUND"].StatusName,
		)
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	userDomainServiceEntity, apiErr := loginUserDomainService.NewLoginUserDomainServiceProps(
		loginUserDomainService.WithLoginUserIdAndEmail(getUser.Email),
		loginUserDomainService.WithLoginUserName(getUser.UserName),
		loginUserDomainService.WithLoginPassword(getUser.Password, loginForm.Password),
	)
	if apiErr != nil {
		log.WithField("apiErr", apiErr).Error("Failed to build login user domain props")
		c.JSON(apiErr.Status, apiErr)
		return loginPresenter, apiErr.Error()
	}

	getUserJson, err := crypto.ConvertStructIntoJson(userDomainServiceEntity)
	if err != nil {
		log.WithError(err).Error("Failed to convert login user domain entity into JSON")
		return loginPresenter, err
	}
	if err := crypto.CopyBeans(getUserJson, &loginPresenter); err != nil {
		log.WithError(err).Error("Failed to copy login user data into presenter")
		return loginPresenter, err
	}

	log.WithField("email", loginPresenter.Email).Info("User logged in successfully")
	return loginPresenter, nil
}
