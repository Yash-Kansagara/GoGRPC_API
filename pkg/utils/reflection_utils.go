package utils

import (
	"errors"
	"reflect"
)

// copies values from one struct to another struct
// it only copies values of same type and same field name
// it does not copy values of unexported fields
// it does not copy values of fields that are not present in the to struct
// it does not copy values of fields that are not present in the from struct
func CopyValues[F any, T any](from *F, to *T) error {
	if from == nil || to == nil {
		return errors.New("from and to cannot be nil")
	}
	fromVal := reflect.ValueOf(from).Elem()
	toVal := reflect.ValueOf(to).Elem()

	if fromVal.Kind() != reflect.Struct {
		return errors.New("from must point to a struct")
	}
	if toVal.Kind() != reflect.Struct {
		return errors.New("to must point to a struct")
	}
	numOfFields := fromVal.NumField()
	for i := 0; i < numOfFields; i++ {
		fromField := fromVal.Field(i)
		toField := toVal.FieldByName(fromVal.Type().Field(i).Name)
		if !toField.IsValid() ||
			!toField.CanSet() ||
			fromField.Type() != toField.Type() {
			continue
		}
		toField.Set(fromField)
	}

	return nil
}
