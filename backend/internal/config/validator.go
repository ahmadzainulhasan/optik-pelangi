package config

import (
	"crypto/subtle"
	"reflect"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/id"
)

func NewValidator(translator ut.Translator) *validator.Validate {
	validate := validator.New()

	_ = id.RegisterDefaultTranslations(validate, translator)

	_ = validate.RegisterValidation("subtle_eqfield", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		currentField, _, _, ok := fl.GetStructFieldOK2()
		if !ok {
			return false
		}
		if field.Kind() != reflect.String || currentField.Kind() != reflect.String {
			return false
		}
		if field.String() == "" || currentField.String() == "" {
			return false
		}
		return subtle.ConstantTimeCompare([]byte(field.String()), []byte(currentField.String())) == 1
	})

	_ = validate.RegisterTranslation("subtle_eqfield", translator, func(ut ut.Translator) error {
		return ut.Add("subtle_eqfield", "{0} harus sama dengan {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), fe.Param())
		return t
	})

	_ = validate.RegisterValidation("empty", func(fl validator.FieldLevel) bool {
		return fl.Field().IsZero()
	})

	_ = validate.RegisterTranslation("empty", translator, func(ut ut.Translator) error {
		return ut.Add("empty", "{0} harus kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	})

	return validate
}
