package d

import "errors"

// 自定义字段名保存全局map
var (
	field = make(map[string]string)
	path = make(map[string]string)
)

const (
	PathUploadDir = "uploadDir"
	PathMenuDir = "menuDir"
)

// 获取字段名
func GetField(fieldKey string) (string, error) {
	val, ok := field[fieldKey]
	if !ok || len(val) == 0 {
		return "", errors.New("field does not exist or is empty")
	}
	return val, nil
}
// 设置字段名
func SetField(key string, value string) {
	field[key] = value
	return
}

// 获取路径
func GetPath(fieldKey string) (string, error) {
	val, ok := path[fieldKey]
	if !ok || len(val) == 0 {
		return "", errors.New("path does not exist or is empty")
	}
	return val, nil
}
// 设置路径
func SetPath(key string, value string) {
	path[key] = value
	return
}