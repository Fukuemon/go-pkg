package validator_test

import (
	"errors"
	"testing"

	"github.com/Fukuemon/go-pkg/validator"
)

// テスト用の構造体
type TestStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Facility string `validate:"ulid"`
	Age      *int   `validate:"omitempty,min=18"`
}

func TestStructValidation(t *testing.T) {
	validator.InitValidator(func(field, tag string) string {
		switch tag {
		case "required":
			return field + " は必須です"
		case "email":
			return field + " は正しいメールアドレス形式である必要があります"
		case "ulid":
			return field + " は有効なULID形式である必要があります"
		default:
			return field + " に無効な値が入力されています"
		}
	})

	tests := []struct {
		name    string
		input   TestStruct
		wantErr error
	}{
		{
			name: "valid data",
			input: TestStruct{
				Name:     "John",
				Email:    "john@example.com",
				Facility: "01AN4Z07BY79KA1307SR9X4MV3", // valid ULID
			},
			wantErr: nil,
		},
		{
			name: "missing name and email",
			input: TestStruct{
				Facility: "01AN4Z07BY79KA1307SR9X4MV3",
			},
			wantErr: errors.New("Name は必須です, Email は必須です"),
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:     "John",
				Email:    "invalid-email",
				Facility: "01AN4Z07BY79KA1307SR9X4MV3",
			},
			wantErr: errors.New("Email は正しいメールアドレス形式である必要があります"),
		},
		{
			name: "invalid ulid",
			input: TestStruct{
				Name:     "John",
				Email:    "john@example.com",
				Facility: "invalid-ulid",
			},
			wantErr: errors.New("Facility は有効なULID形式である必要があります"),
		},
		{
			name: "age below minimum",
			input: TestStruct{
				Name:     "John",
				Email:    "john@example.com",
				Facility: "01AN4Z07BY79KA1307SR9X4MV3",
				Age:      intPointer(17),
			},
			wantErr: errors.New("Age に無効な値が入力されています"),
		},
	}

	// テストの実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.StructValidation(tt.input)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("expected no error but got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Fatalf("expected error %v but got no error", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
			}
		})
	}
}

// intのポインタを作成するためのヘルパー関数
func intPointer(i int) *int {
	return &i
}
