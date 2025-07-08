package handler

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// handlerにあると便利なヘルパー関数群 //
// リクエストに必要な型とデータを渡すことで空なデータがないかチェックするヘルパー関数
// validate struct
func ValidateStruct[T any](s *T) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("input must be a pointer struct")
	}

	elem := v.Elem()
	t := elem.Type()
	uuidType := reflect.TypeOf(uuid.UUID{})

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := field.Type()
		fieldName := t.Field(i).Name

		// bool型は常に値がセットされているとみなす
		if fieldType.Kind() == reflect.Bool {
			continue
		}

		if field.IsZero() {
			return fmt.Errorf("field %s is required and cannnot be empty", fieldName)
		}

		switch fieldType {
		case uuidType:
			id, ok := field.Interface().(uuid.UUID)
			if !ok {
				return fmt.Errorf("internal error: could not cast field '%s' to uuid.UUID", fieldName)
			}
			_, err := uuid.Parse(id.String())
			if err != nil {
				return fmt.Errorf("field '%s' contains a malformed UUID: %v", fieldName, err)
			}
		}
	}

	return nil
}
