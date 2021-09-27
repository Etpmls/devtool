package d_recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var Recaptcha recaptcha

type recaptcha struct {
	Secret string // recaptcha密钥
	Optional optionalRecaptcha
	enable bool
}

type optionalRecaptcha struct {
	Host string	// recaptcha默认域名，www.google.com，如中国大陆，请设置www.recaptcha.net
}

func (this *recaptcha) Init() error {
	if this.Secret == "" {
		return errors.New("secret is required")
	}

	// recaptcha默认域名，www.google.com，如中国大陆，请设置www.recaptcha.net
	if this.Optional.Host == "" {
		this.Optional.Host = "www.google.com"
	}

	this.enable = true
	return nil
}

// 获取启动的状态
func (this *recaptcha) GetEnabledStatus() bool {
	return this.enable
}

// recaptcha返回格式
type VerifyCaptchaResponse struct {
	Success bool `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}
func (this *recaptcha) VerifyCaptcha(response string) (error) {
	// 与谷歌通讯
	resp, err := http.PostForm("https://" + this.Optional.Host + "/recaptcha/api/siteverify", url.Values{"secret":{this.Secret}, "response":{response}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res VerifyCaptchaResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if res.Success == false {
		return errors.New(strings.Join(res.ErrorCodes, ","))
	}

	return nil
}