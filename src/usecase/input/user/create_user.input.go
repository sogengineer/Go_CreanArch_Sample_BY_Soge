package inputUser

import (
	"regexp"

	"github.com/Go_CleanArch/common/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

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
