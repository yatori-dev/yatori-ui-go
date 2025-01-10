package utils

import "reflect"

func PageFormat(obj interface{}) {
	// 反射获取字段page和size的值，如果为空则设置默认值
	if obj == nil {
		return
	}
	reflectValue := reflect.ValueOf(obj)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return
	}
	reflectType := reflectValue.Type()
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		if fieldType.Name == "Page" && field.Kind() == reflect.Int {
			if field.Int() == 0 {
				field.SetInt(1)
			}
		}
		if fieldType.Name == "Size" && field.Kind() == reflect.Int {
			if field.Int() == 0 {
				field.SetInt(10)
			}
		}
	}
}
