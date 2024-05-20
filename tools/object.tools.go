package tools

import "reflect"

func MapFields(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()
	srcType := srcVal.Type()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcType.Field(i).Name)

		if dstField.IsValid() && dstField.Type() == srcField.Type() {
			dstField.Set(srcField)
		}
	}
}
