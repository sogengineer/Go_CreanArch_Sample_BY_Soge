# Go による Clean Architecture の原則を参考に実装したサンプル構成
このアプリケーションは、Goにおいて Clean Architecture の原則を参考に実装したサンプルアプリケーションです。アプリケーションは階層化されたアーキテクチャに従っており、関心事の分離、モジュール性、テスト容易性、および保守性を促進しています。

## アーキテクチャの概要
アプリケーションは以下の層で構成されています。
- インターフェース層（Interface Adapter Layer）
  - src/interface_adapter配下
- アプリケーション層（Application Layer）
  - src/usecase配下
- ドメイン層（Domain Layer）
  - src/domain配下
- インフラストラクチャ層（Infrastructure Layer）
  - src/infrastructure配下

### インターフェース層（Interface Adapter Layer）
インターフェース層は、外部からの入力を受け取り、適切なユースケースを呼び出すための層。この層では、以下のような要素が含まれる。
#### コントローラー（Controller）
コントローラーは、外部からのリクエストを受け取り、適切なユースケースを呼び出す役割を持つ

例：
```golang
type UserController struct {
	userService userService.UserService
}

// NewUserController is the constructor for UserController
func NewUserController(userService userService.UserService) *UserController {
	return &UserController{userService: userService}
}

// Create action: POST /users
func (uc *UserController) UserController(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := uc.userService.CreateUserService(ctx, c)
	if err != nil {
		fmt.Println(err)
	} else {
		c.JSON(
			status.SuccessStatusMap["CREATED"].StatusCode,
			result,
		)
	}
}
```
この関数は、以下のようなクリーンアーキテクチャの原則に従っています。

- 外部からのリクエストを受け取り、適切なユースケースを呼び出す。
- ユースケースからの結果をレスポンスとして返す。
- ユースケースの実装の詳細に関与しない。

#### ゲートウェイ（Gateway）
外部システムとのインタラクションを抽象化する役割を持つ要素
例：
```golang
func (userRepository *userRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    err := dbConnect.Find(ctx, "email = ?", []interface{}{email}, &user)
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("条件に一致するレコードが見つかりません: %w", err)
    } else if err != nil {
        return nil, fmt.Errorf("DB検索に失敗しました: %w", err)
    }
    return &user, nil
}
```
■特徴
- データベースアクセスを抽象化し、提供先(アプリケーション層)へ詳細を隠蔽している。
- 同階層のエンティティを使用してデータの読み書きを行っている。
  - gateway内のエンティティはテーブル定義と同様にしている
- インフラストラクチャ層の実装（dbConnect）を使用している。※後ほど説明
  - 汎用的なDB操作の詳細はインフラストラクチャ層で実装し、特定のスキーマおよびテーブルの指定、操作はゲートウェイで実装している。(プロジェクトや場合によりますが、自分はその方法を採用しています。)

<br><br>

### UseCase層（Application Layer）
ユースケースの実装を担当する層。
#### インプット（Input）
リクエストパラメータで処理の対象となる項目を表すデータを定義する場所
例：
```golang
type CreateUserForm struct {
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	CreatedFlg bool   `json:"createdFlg"`
}
// CreateUserForm専用入力バリデーション
func (createUserForm CreateUserForm) CreateUserValidate() []errors.ApiErrMessage {
	var apiErrMessages []errors.ApiErrMessage
	createUserFormValidation := validation.ValidateStruct(&createUserForm,
		validation.Field(
			&createUserForm.UserName,
			validation.Required.Error("ユーザー名を入力してください"),
			validation.Length(1, 30).Error("ユーザー名は 30文字以内で入力してください"),
		),
		validation.Field(
			&createUserForm.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			is.Email.Error("正しいメールアドレスを入力してください"),
			validation.RuneLength(5, 40).Error("メールアドレスは 5～40文字です"),
		),
		validation.Field(
			&createUserForm.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.Length(8, 16).Error("パスワードは8〜16桁で入力してください"),
			validation.Match(regexp.MustCompile("^*[a-z].*$")).Error("パスワードは半角の英大文字、英小文字、数字を含む形式にしてください"),
			validation.Match(regexp.MustCompile("^*[A-Z].*$")).Error("パスワードは半角の英大文字、英小文字、数字を含む形式にしてください"),
			validation.Match(regexp.MustCompile("^*[0-9].*$")).Error("パスワードは半角の英大文字、英小文字、数字を含む形式にしてください"),
		),
	)
	if err := createUserFormValidation; err != nil {
		errors.AddValidationErrors(&apiErrMessages, err, nil)
		return apiErrMessages
	}
	return nil
}
```
■特徴
- 各フィールドには、JSONタグを使用しAPIリクエストのJSONキーとの対応を定義
- CreateUserValidateメソッドは、ozzo-validationライブラリを使用して入力データのバリデーションを行っている
  - バリデーションエラーがある場合は、本サンプル実装独自の型であるApiErrMessageスライスを返す

<br>

#### アウトプット（Output）
出力データを定義する箇所

<br>

#### サービス（Service）
```golang
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
```
■処理の流れ
- 1.リクエストボディからJSONデータをCreateUserForm構造体にバインド
- 2.CreateUserValidateメソッドを呼び出し、入力データのバリデーションを行います。バリデーションエラーがある場合は、エラーレスポンスを返す
- 3.UserRepositoryを使用して、登録済みのメールアドレスを再登録しようとしていないかチェック
- 4.createUserDomain.NewCreateCommonUserFactoryProps関数を呼び出して、登録するユーザー情報のビルドを行う
- 5.ビルドしたユーザー情報をJSONに変換し、UserRepositoryのCreateUserメソッドを呼び出してユーザー登録を行う
- 6.登録されたユーザー情報をCreateUserPresenter構造体に変換し、レスポンスとして返す

■Serviceの役割
- 受け取ったInputのパラメータをもとに、集約を跨ったりコアビジネスロジックを経由したりして、定義されているOutputパラメータへとデータが整形されていく過程を司る箇所
- ドメイン層とインターフェース層の仲介役を担っている

<br>

#### リポジトリインターフェース(Repository Interface)
ユースケースから使用するリポジトリのインターフェースを定義
```golang
type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, userJson []byte) (*entity.User, error)
}
```
アプリケーション層では、このリポジトリインターフェースを使用してデータアクセスを行う。実際のリポジトリの実装は、インターフェース層・ゲートウェイで行われる。

<br>

**※なぜ、リポジトリのインターフェースは提供側でなく利用側の階層に定義するの？**
- 利用側でインターフェースを定義することで、利用側のコードがリポジトリの具体的な実装ではなくインターフェースに依存するようになり依存関係の方向を逆転させることができるから。(依存性逆転の法則(DI)と言います)
  - DIにより利用側のコードが実装の詳細(Gateway層/Repository)から切り離され柔軟性が高くなる他、リポジトリ層をモック化することが容易となり、テスタビリティが向上する効果がある
- ただし、これはGoにおける一般的な傾向であり、プロジェクトの要件や設計方針によっては、提供側でインターフェースを定義することもあります。
  - 重要なのは**コードの柔軟性、保守性、テスト容易性を高めること**

<br>
<br>

### Domain層
アプリケーションのコアビジネスロジック、エンティティ、およびアプリケーションの契約を定義する箇所
#### ドメインサービス(Domain Service)
複数のエンティティやバリューオブジェクトを調整し、ビジネスロジックを実装するオブジェクトのこと
