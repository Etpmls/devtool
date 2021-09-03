package d

import (
	"errors"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	v "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)
// 主要变量
var (
	Validator validator
	Validate *v.Validate
)
// 参数配置
type validator struct {
	Optional optionalValidator
}
// 可选参数配置
type optionalValidator struct {
	Fallback string
	SupportedLocales struct {
		En bool
		Zh bool
	}
}
// 定义的变量默认值
var ValidatorDefine = validatorDefine{
	En: "en",
	Zh: "zh",
}
// 定义的变量
type validatorDefine struct {
	En string
	Zh string
	TransEn ut.Translator
	TransZh ut.Translator
	// use a single instance , it caches struct info
	uni      *ut.UniversalTranslator
}

func (this *validator) Init() error {
	Validate = v.New()

	var (
		fallback locales.Translator
		supportedLocales []locales.Translator
	)

	switch strings.ToLower(this.Optional.Fallback) {
	case ValidatorDefine.En:
		fallback = en.New()
	case ValidatorDefine.Zh:
		fallback = zh.New()
	case "":
	default:
		return errors.New("Fallback is not configured or the language " + this.Optional.Fallback + " is temporarily not supported.")
	}

	if this.Optional.SupportedLocales.En {
		supportedLocales = append(supportedLocales, en.New())
	}

	if this.Optional.SupportedLocales.Zh {
		supportedLocales = append(supportedLocales, zh.New())
	}

	// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go
	ValidatorDefine.uni = ut.New(fallback, supportedLocales...)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	if this.Optional.SupportedLocales.En {
		ValidatorDefine.TransEn, _ = ValidatorDefine.uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(Validate, ValidatorDefine.TransEn)
	}
	if this.Optional.SupportedLocales.Zh {
		ValidatorDefine.TransZh, _ = ValidatorDefine.uni.GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(Validate, ValidatorDefine.TransZh)
	}

	return nil
}

// 翻译到英文
func (this *validator) TranslateEn(err error) error {
	if !this.Optional.SupportedLocales.En {
		return errors.New("en language is not enabled.")
	}

	errs := err.(v.ValidationErrors)
	var s []string
	for _, e := range errs {
		// can translate each error one at a time.
		s = append(s, e.Translate(ValidatorDefine.TransEn))
	}
	return errors.New(strings.Join(s, ","))
}

// 翻译到中文
func (this *validator) TranslateZh(err error) error {
	if !this.Optional.SupportedLocales.Zh {
		return errors.New("zh language is not enabled.")
	}

	errs := err.(v.ValidationErrors)
	var s []string
	for _, e := range errs {
		// can translate each error one at a time.
		s = append(s, e.Translate(ValidatorDefine.TransZh))
	}
	return errors.New(strings.Join(s, ","))
}