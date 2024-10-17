package validator_test

import (
	"errors"
	"testing"

	"github.com/Fukuemon/go-pkg/ulid"
	validator "github.com/Fukuemon/go-pkg/validator/gin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// テストの初期化
func init() {
	validator.InitTagErrorMessages(map[string]string{
		"required": "カスタム: このフィールドは必須です",
		"ulid":     "カスタム: 無効なULID形式です",
		"int":      "カスタム: 無効な整数形式です",
	})
}

// Param関数のテスト
func TestParam(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		rules      []string
		wantErr    error
	}{
		{
			name:       "Valid ULID",
			paramName:  "facility_id",
			paramValue: ulid.NewULID(),
			rules:      []string{"required", "ulid"},
			wantErr:    nil,
		},
		{
			name:       "Missing required field",
			paramName:  "facility_id",
			paramValue: "",
			rules:      []string{"required"},
			wantErr:    errors.New("カスタム: このフィールドは必須です"),
		},
		{
			name:       "Invalid ULID",
			paramName:  "facility_id",
			paramValue: "invalid-ulid",
			rules:      []string{"required", "ulid"},
			wantErr:    errors.New("カスタム: 無効なULID形式です"),
		},
		{
			name:       "Valid int",
			paramName:  "department_id",
			paramValue: "123",
			rules:      []string{"required", "int"},
			wantErr:    nil,
		},
		{
			name:       "Invalid int",
			paramName:  "department_id",
			paramValue: "invalid-int",
			rules:      []string{"required", "int"},
			wantErr:    errors.New("カスタム: 無効な整数形式です"),
		},
	}

	// gin.Contextのモックを作成
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ginのモックコンテキストを作成
			ctx, _ := gin.CreateTestContext(nil)
			ctx.Params = gin.Params{
				gin.Param{Key: tt.paramName, Value: tt.paramValue},
			}

			// Param関数をテスト
			result := validator.Param(ctx, tt.paramName, tt.rules...)

			// エラーチェック
			if tt.wantErr != nil {
				assert.EqualError(t, result.Err, tt.wantErr.Error())
			} else {
				assert.NoError(t, result.Err)
			}
		})
	}
}

// ParamsValidation関数のテスト
func TestParamsValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  []validator.ParamValidation
		wantErr error
	}{
		{
			name: "All valid",
			params: []validator.ParamValidation{
				{
					ParamName:  "facility_id",
					ParamValue: ulid.NewULID(),
					Rules:      []string{"required", "ulid"},
					Err:        nil,
				},
				{
					ParamName:  "department_id",
					ParamValue: "123",
					Rules:      []string{"required", "int"},
					Err:        nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "One invalid ULID",
			params: []validator.ParamValidation{
				{
					ParamName:  "facility_id",
					ParamValue: "invalid-ulid",
					Rules:      []string{"required", "ulid"},
					Err:        errors.New("カスタム: 無効なULID形式です"),
				},
				{
					ParamName:  "department_id",
					ParamValue: "123",
					Rules:      []string{"required", "int"},
					Err:        nil,
				},
			},
			wantErr: errors.New("カスタム: 無効なULID形式です"),
		},
		{
			name: "One invalid int",
			params: []validator.ParamValidation{
				{
					ParamName:  "facility_id",
					ParamValue: ulid.NewULID(),
					Rules:      []string{"required", "ulid"},
					Err:        nil,
				},
				{
					ParamName:  "department_id",
					ParamValue: "invalid-int",
					Rules:      []string{"required", "int"},
					Err:        errors.New("カスタム: 無効な整数形式です"),
				},
			},
			wantErr: errors.New("カスタム: 無効な整数形式です"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ParamsValidation関数をテスト
			err := validator.ParamsValidation(tt.params...)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
