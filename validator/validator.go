package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate            *validator.Validate
	customErrorMessages func(field, tag string) string // カスタムエラーメッセージ生成関数
)

// InitValidator はバリデーションの初期化を行い、カスタムメッセージ関数も登録できる
func InitValidator(customMsgFunc func(field, tag string) string) {
	if validate == nil {
		validate = validator.New()
	}

	// カスタムメッセージ関数をセット
	customErrorMessages = customMsgFunc
}

// 構造体（Body）のバリデーションを実行する
func StructValidation(s interface{}) error {
	if validate == nil {
		return errors.New("validator is not initialized")
	}

	// バリデーションの実行
	err := validate.Struct(s)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return ValidationError(ve)
		}
		return err
	}
	return nil
}

// 構造体バリデーション時のエラーメッセージを生成する
func ValidationError(ve validator.ValidationErrors) error {
	errorMessages := make([]string, len(ve))
	for i, fieldErr := range ve {
		// カスタムメッセージ関数が定義されている場合はそれを使う
		if customErrorMessages != nil {
			errorMessages[i] = customErrorMessages(fieldErr.Field(), fieldErr.Tag())
		} else {
			// デフォルトメッセージを使用
			errorMessages[i] = getDefaultErrorMessage(fieldErr)
		}
	}
	return errors.New(strings.Join(errorMessages, ", "))
}

func getDefaultErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " は必須です"
	case "email":
		return fe.Field() + " は正しいメールアドレス形式である必要があります"
	case "ulid":
		return fe.Field() + " は有効なULID形式である必要があります"
	default:
		return fe.Field() + " に無効な値が入力されています"
	}
}
