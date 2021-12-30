package d_gin

import (
	d "github.com/etpmls/devtool"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

// 上传图片
func ImageUpload(c *gin.Context, subDir string) (filePath string, err error) {
	// 如果subDir未设置，默认为日期
	if subDir == "" {
		subDir = time.Now().Format("20060102")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}

	extension, err := d.ImageValidate(f)
	if err != nil {
		return "", err
	}

	// Make Dir
	path := d.GetUploadPath() + subDir + "/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	// UUID File name
	u := strings.ReplaceAll(uuid.New().String(), "-", "")

	file_path := path + u + "." + extension
	err = c.SaveUploadedFile(file, file_path)
	if err != nil {
		return "", err
	}

	return file_path, nil
}

// Delete Image
// 删除图片
type ImageDeleteRequest struct {
	Path string	`json:"path" binding:"required"`
}
func ImageDelete(c *gin.Context, translator ut.Translator) error {
	// 校验请求
	var json ImageDeleteRequest
	err := Validate(c, &json, translator)
	if err != nil {
		return err
	}

	// 校验操作的路径是否合法
	err = d.FilePathValidate(json.Path, []string{d.GetUploadPath()})
	if err != nil {
		return err
	}

	// 文件删除
	err = d.FileDelete(json.Path)
	if err != nil {
		return err
	}

	return nil
}

