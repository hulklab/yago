package validatelib

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
)

func New(trans ut.Translator) *validator.Validate {
	v := validator.New()
	if trans == nil {
		return v
	}

	lang := trans.Locale()
	switch lang {
	case "zh":
		_ = zh.RegisterDefaultTranslations(v, trans)
	case "en":
		_ = en.RegisterDefaultTranslations(v, trans)
	default:
		_ = zh.RegisterDefaultTranslations(v, trans)
	}

	return v

}
