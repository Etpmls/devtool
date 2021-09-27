package d

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Config config
)

type config struct {
	FilePath string
	ConfigAddr interface{}
	enable bool
}

func (this *config) Init() error {
	// If it is empty, skip initialization
	// 如果为空，则跳过初始化
	if this.FilePath == "" {
		return errors.New("config file path is empty")
	}

	if this.ConfigAddr == nil {
		return errors.New("config address is nil")
	}

	// 读取文件
	b, err := ioutil.ReadFile(this.FilePath)
	if err != nil {
		return err
	}

	// 把值解析到结构体
	err = yaml.Unmarshal(b, this.ConfigAddr)
	if err != nil {
		return err
	}

	this.enable = true
	return nil
}

// 获取启动的状态
func (this *config) GetEnabledStatus() bool {
	return this.enable
}