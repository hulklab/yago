package validatelib

import (
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
)

var v *validator.Validate
var once sync.Once

func Ins(trans ut.Translator) *validator.Validate {
	once.Do(func() {
		v = validator.New()
	})

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

	registerPhoneValidator(v)

	return v
}
