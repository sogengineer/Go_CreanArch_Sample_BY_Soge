package errors

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApiErr struct {
	Messages []ApiErrMessage `json:"messages"`
	Status   int             `json:"status"`
	Detail   string          `json:"detail"`
}

type ApiErrMessage struct {
	Key           string         `json:"key"`
	Value         string         `json:"value"`
	ErrorLocation *ErrorLocation `json:"errorLocation,omitempty"` // Optional field
}

type ErrorLocation struct {
	Object        string         `json:"object"`
	Index         *int           `json:"index,omitempty"`         // Optional field
	ErrorLocation *ErrorLocation `json:"errorLocation,omitempty"` // Optional field
}

func OutputApiError(messages []ApiErrMessage, state int, detail string) *ApiErr {
	var apiErr ApiErr
	apiErr.Messages = messages
	apiErr.Status = state
	apiErr.Detail = detail

	return &apiErr
}

// addValidationErrors はバリデーションエラーを受け取り、ApiErrMessage スライスに追加します。
func AddValidationErrors(apiErrMessages *[]ApiErrMessage, err error, errorLocation *ErrorLocation) {
	if validationErrors, ok := err.(validation.Errors); ok {
		for field, e := range validationErrors {
			*apiErrMessages = append(*apiErrMessages, ApiErrMessage{
				Key:           field,     // エラーの発生したフィールド名
				Value:         e.Error(), // エラーメッセージ
				ErrorLocation: errorLocation,
			})
		}
	}
}

func (apiErr ApiErr) Error() error {
	var errMessages []string
	for _, msg := range apiErr.Messages {
		errMsg := fmt.Sprintf("	key: %s, value: %s", msg.Key, msg.Value)
		if msg.ErrorLocation != nil {
			errMsg += fmt.Sprintf(" (object: %s", msg.ErrorLocation.Object)
			if msg.ErrorLocation.Index != nil {
				errMsg += fmt.Sprintf(", index: %d", *msg.ErrorLocation.Index)
			}
			if msg.ErrorLocation.ErrorLocation != nil {
				errMsg += fmt.Sprintf(", errorLocation: %+v", *msg.ErrorLocation.ErrorLocation)
			}
			errMsg += ")"
		}
		errMessages = append(errMessages, errMsg)
	}
	errMsg := strings.Join(errMessages, ", ")
	return fmt.Errorf(`	{
		"messages": [
		%s
		],
		"status": %d,
		"detail": "%s"
	}`, errMsg, apiErr.Status, apiErr.Detail)
}
