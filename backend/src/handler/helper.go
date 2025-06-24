package handler

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// handlerにあると便利なヘルパー関数群 //
// リクエストに必要な型とデータを渡すことで空なデータがないかチェックするヘルパー関数
func ValidateStruct[T any](s *T) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("input must be a pointer struct")
	}

	// ポインタから構造体を取得
	elem := v.Elem()
	t := elem.Type()

	// 詳細なチェックを行う型の定義
	uuidType := reflect.TypeOf(uuid.UUID{})

	// 構造体のフィールドをループ
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldName := t.Field(i).Name

		if field.IsZero() {
			return fmt.Errorf("field %s is required and cannnot be empty", fieldName)
		}

		switch field.Type() {
		case uuidType:
			// 型がuuid.UUIDの場合の処理
			id, ok := field.Interface().(uuid.UUID)
			if !ok {
				return fmt.Errorf("internal error: could not cast field '%s' to uuid.UUID", fieldName)
			}
			_, err := uuid.Parse(id.String())
			if err != nil {
				return fmt.Errorf("field '%s' contains a malformed UUID: %v", fieldName, err)
			}
		default:
		}
	}

	return nil
}
