package config

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/id"
)

func NewValidator(translator ut.Translator) *validator.Validate {
	validate := validator.New()

	_ = id.RegisterDefaultTranslations(validate, translator)

	return validate
}
