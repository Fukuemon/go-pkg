package validator

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Fukuemon/go-pkg/ulid"
	"github.com/gin-gonic/gin"
)

// タグごとのエラーメッセージを保持するマップ
var tagErrorMessages map[string]string

func InitTagErrorMessages(customMessages map[string]string) {
	// タグごとのデフォルトエラーメッセージ
	tagErrorMessages = map[string]string{
		"required": "{field}は必須です",
		"ulid":     "{field}は無効なULID形式です",
		"int":      "{field}は無効な整数形式です",
	}

	// カスタムメッセージを上書き
	for tag, message := range customMessages {
		if message != "" {
			tagErrorMessages[tag] = message
		}
	}
}

// ParamValidation 構造体はバリデーション結果を保持
type ParamValidation struct {
	ParamName  string
	ParamValue string
	Rules      []string
	Err        error
}

// 単一のパスパラメータに対するバリデーション
func Param(ctx *gin.Context, paramName string, rules ...string) ParamValidation {
	paramValue := ctx.Param(paramName)

	// バリデーション実行
	var validationErr error
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if err := PathParamSingleValidation(paramName, paramValue, rule); err != nil {
			validationErr = err
			break
		}
	}

	// ParamValidation構造体に結果を格納
	return ParamValidation{
		ParamName:  paramName,
		ParamValue: paramValue,
		Rules:      rules,
		Err:        validationErr,
	}
}

// 複数のパスパラメータに対するバリデーション
func ParamsValidation(params ...ParamValidation) error {
	var validationErrors []string

	for _, param := range params {
		if param.Err != nil {
			validationErrors = append(validationErrors, param.Err.Error())
		}
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, ", "))
	}

	return nil
}

// PathParamSingleValidation は単一のパラメータバリデーション
func PathParamSingleValidation(paramName, paramValue, tag string) error {
	switch tag {
	case "required":
		if paramValue == "" {
			return errors.New(getTagErrorMessage(paramName, tag))
		}
	case "ulid":
		if !ulid.IsValid(paramValue) {
			return errors.New(getTagErrorMessage(paramName, tag))
		}
	case "int":
		if _, err := strconv.Atoi(paramValue); err != nil {
			return errors.New(getTagErrorMessage(paramName, tag))
		}
	default:
		return errors.New(getTagErrorMessage(paramName, tag))
	}
	return nil
}

// 各バリデーションタグに対するエラーメッセージを返す
func getTagErrorMessage(paramName, tag string) string {
	message := tagErrorMessages[tag]
	return strings.ReplaceAll(message, "{field}", paramName)
}
