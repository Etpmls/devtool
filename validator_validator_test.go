package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestValidator(t *testing.T) {
	err := d.Validator.Init()
	if err != nil {
		t.Fatal(err)
	}
	translateAll()
	translateIndividual()

	// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go#L12
	d.Validator.Optional.Fallback = d.ValidatorDefine.En
	d.Validator.Optional.SupportedLocales.En = true
	d.Validator.Optional.SupportedLocales.Zh = true
	err = d.Validator.Init()
	if err != nil {
		t.Fatal(err)
	}

	translateAll()
	translateIndividual()
	translateOverride()
}

func translateAll() {

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
		fmt.Println(d.Validator.TranslateEn(err))
		fmt.Println(d.Validator.TranslateZh(err))
	}
}

func translateIndividual() {

	type User struct {
		Username string `validate:"required"`
	}

	var user User

	err := d.Validate.Struct(user)
	if err != nil {
		fmt.Println(d.Validator.TranslateEn(err))
		fmt.Println(d.Validator.TranslateZh(err))
	}
}

func translateOverride() {

	d.Validate.RegisterTranslation("required", d.ValidatorDefine.TransEn, func(ut ut.Translator) error {
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
		fmt.Println(d.Validator.TranslateEn(err))
		fmt.Println(d.Validator.TranslateZh(err))
	}
}