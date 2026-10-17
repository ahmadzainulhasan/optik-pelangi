package config

import (
	"github.com/go-playground/locales/id"
	"github.com/go-playground/universal-translator"
)

func NewTranslator() ut.Translator {
	indonesianLocale := id.New()
	universal := ut.New(indonesianLocale, indonesianLocale)
	translator, _ := universal.GetTranslator("id")

	return translator
}
