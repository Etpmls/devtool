package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"testing"
)

func TestValidator(t *testing.T) {
	// 初始化，不添加i18n的simple案例
	// https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
	err := d.Validator.Init()
	if err != nil {
		t.Fatal(err)
	}
	translateAll(nil)
	translateIndividual(nil)

	var (
		uni      *ut.UniversalTranslator
		transEn	ut.Translator
		transZh	ut.Translator
	)

	// 增加I18n
	// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go#L12
	d.Validator.Optional.OverrideInit = func() error {
		// NOTE: ommitting allot of error checking for brevity
		en := en.New()
		zh := zh.New()
		uni = ut.New(en, en, zh)
		transEn, _ = uni.GetTranslator("en")
		transZh, _ = uni.GetTranslator("zh")
		en_translations.RegisterDefaultTranslations(d.Validate, transEn)
		zh_translations.RegisterDefaultTranslations(d.Validate, transZh)
		return nil
	}

	err = d.Validator.Init()
	if err != nil {
		t.Fatal(err)
	}

	translateAll(transEn)
	translateIndividual(transZh)
	translateOverride(transZh)
}

func translateAll(trans ut.Translator) {

	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 string `validate:"required,gt=1"`
	}

	user := User{
		Username: "Joeybloggs",
		Tagline:  "This tagline is way too long.",
		Tagline2: "1",
	}

	err := d.Validate.Struct(user)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'

		if trans != nil {
			fmt.Println(errs.Translate(trans))
		} else {
			fmt.Println(errs)
		}
	}
}

func translateIndividual(trans ut.Translator) {

	type User struct {
		Username string `validate:"required"`
	}

	var user User

	err := d.Validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			if trans != nil {
				fmt.Println(e.Translate(trans))
			} else {
				fmt.Println(errs)
			}
		}
	}
}

func translateOverride(trans ut.Translator) {

	d.Validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	type User struct {
		Username string `validate:"required"`
	}

	var user User

	err := d.Validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			fmt.Println(e.Translate(trans))
		}
	}
}