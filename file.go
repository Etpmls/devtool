package d

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
)

// 获取上传路径，可自定义，如无设置，则获取默认路径
func GetUploadPath() string {
	// 如果没有自定义路径，将使用默认路径
	uploadDir, err := GetPath(PathUploadDir)
	if err != nil {
		uploadDir = "storage/upload/"
	}
	return uploadDir
}

// 获取上传路径，可自定义，如无设置，则获取默认路径
func SetUploadPath(path string)  {
	// 如果没有自定义路径，将使用默认路径
	SetPath(PathUploadDir, path)
	return
}

// 验证文件是否为图片，如果是返回文件后缀名
func ImageValidate(f multipart.File) (extension string, err error) {
	// 识别图片类型
	_, image_type, _ := image.Decode(f)

	// 获取图片的类型
	switch image_type {
	case `jpeg`:
		return "jpeg", nil
	case `png`:
		return "png", nil
	case `gif`:
		return "git", nil
	case `bmp`:
		return "bmp", nil
	default:
		return "", errors.New("This is not an image file, or the image file format is not supported!")
	}
}

// 文件校验，是否在允许操作的目录中
func FilePathValidate(path string, allowDir []string) error {
	// 把./test/前缀统一转化格式为test/，便于匹配是否在允许的上传路径内
	reg := regexp.MustCompile("^(\\.?)+/")
	path = reg.ReplaceAllString(path, "")

	for _, v := range allowDir {
		v = reg.ReplaceAllString(v, "")
		// 验证是否在上传目录路径内，如果不是，非法删除
		if len(path) > len(v) && strings.Contains(path[:len(v)], v) {
			return  nil
		}
	}

	return  errors.New("illegal request path")
}

// 验证文件状态，以及是否为目录，如果文件不存在则返回nil
func FileCheck(filePath string) error {
	f, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// 如果文件是目录
	if f.IsDir() {
		return errors.New("the path is a directory")
	}
	return nil
}

// 删除文件
func FileDelete(path string) error {
	// 检查是否存在,或者是否为一个目录
	err := FileCheck(path)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// Batch delete any type of files
// 批量删除任何类型文件
func FileBatchDelete(path []string) (err error) {
	// 第一遍遍历先验证
	for _, v := range path {
		// Validate If a File
		err = FileCheck(v)
		if err != nil {
			return err
		}
	}
	// 第二遍遍历再删除
	for _, v := range path {
		// Delete Image
		_ = os.Remove(v)
	}

	return err
}