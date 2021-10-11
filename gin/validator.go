package d_gin

import (
	"errors"
	d "github.com/Etpmls/devtool"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

// 通用的校验方法
func Validate(c *gin.Context, structAddr interface{}, translator ut.Translator) error {
	// 绑定数据
	if err := c.ShouldBindJSON(structAddr); err != nil {
		return err
	}

	// 验证数据
	err := d.ValidatorClient.Struct(structAddr)
	if err != nil {
		if translator == nil {
			return err
		}

		errs := err.(validator.ValidationErrors)
		var s []string
		for _, e := range errs {
			// can translate each error one at a time.
			s = append(s, e.Translate(translator))
		}
		return errors.New(strings.Join(s, ","))
	}

	return nil
}
