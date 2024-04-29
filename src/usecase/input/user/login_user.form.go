package inputUser

import (

	"github.com/Go_CleanArch/common/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

// ログイン
type LoginForm struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

// loginForm専用入力バリデーション
func (loginForm LoginForm) LoginValidate() []errors.ApiErrMessage {
	var apiErrMessages []errors.ApiErrMessage
	loginFormValidation := validation.ValidateStruct(&loginForm,
		validation.Field(
			&loginForm.Email,
			validation.Required.Error("正しいメールアドレスを入力してください"),
			is.Email.Error("正しいメールアドレスを入力してください"),
			validation.RuneLength(5, 40).Error("メールアドレスは 5～40 文字です"),
		),
		validation.Field(
			&loginForm.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.Length(8, 16).Error("パスワードは8〜16桁で入力してください"),
		),
	)
	if err := loginFormValidation; err != nil {
		errors.AddValidationErrors(&apiErrMessages, err, nil)
		return apiErrMessages
	}
	return nil
}