package d

import (
	"reflect"
)

func CopyStructValue(source interface{}, target interface{})  {
	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(target)

	// 如果target传过来的不是指针，往下进行set操作会报错，原样返回
	if targetVal.Kind() != reflect.Ptr {
		return
	}

	// 如果source传过来的是指针，会报错：reflect: call of reflect.Value.NumField on ptr Value [recovered]，从指针取值
	if sourceVal.Kind() == reflect.Ptr {
		sourceVal = sourceVal.Elem()
	}

	// 从指针中取target的值
	targetVal = targetVal.Elem()

	for i := 0; i < sourceVal.NumField(); i++ {
		name := sourceVal.Type().Field(i).Name
		ok := targetVal.FieldByName(name).IsValid();
		if ok && targetVal.FieldByName(name).Kind() == sourceVal.Field(i).Kind() {
			targetVal.FieldByName(name).Set(sourceVal.Field(i))
		}
	}
	return
}
