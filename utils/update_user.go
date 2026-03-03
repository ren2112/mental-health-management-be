package utils

import "reflect"

// 自动生成更新字段（忽略零值）
func BuildUpdateMap(req interface{}) map[string]interface{} {

	updateData := make(map[string]interface{})

	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		typeField := typ.Field(i)

		jsonTag := typeField.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 指针类型（判断是否传值）
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				updateData[jsonTag] = field.Elem().Interface()
			}
			continue
		}

		// 普通类型过滤零值
		zero := reflect.Zero(field.Type()).Interface()
		if !reflect.DeepEqual(field.Interface(), zero) {
			updateData[jsonTag] = field.Interface()
		}
	}

	return updateData
}
