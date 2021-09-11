package d

import (
	v "github.com/go-playground/validator/v10"
)
// 主要变量
var (
	Validator       validator
	ValidatorClient *v.Validate
)
// 参数配置
type validator struct {
	Optional validatorOptional
	enable bool
}
// 可选参数配置
type validatorOptional struct {
	OverrideInit func() error
}

func (this *validator) Init() error {
	ValidatorClient = v.New()
	if this.Optional.OverrideInit != nil {
		// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go
		err := this.Optional.OverrideInit()
		if err != nil {
			return err
		}
	}
	this.enable = true
	return nil
}

// 获取启动的状态
func (this *validator) GetEnabledStatus() bool {
	return this.enable
}