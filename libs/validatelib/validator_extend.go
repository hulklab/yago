package validatelib

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hulklab/yago"
)

func registerPhoneValidator(v *validator.Validate) {
	// 注册自定义验证器和翻译
	// 例：添加一个手机号验证器
	// 添加手机号验证器
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
		reg := regexp.MustCompile(regular)
		return reg.MatchString(fl.Field().String())
	})

	// 添加手机号验证器的翻译
	_ = v.RegisterTranslation("phone", yago.GetTranslator(), func(ut ut.Translator) error {
		var e error
		switch ut.Locale() {
		case "zh":
			e = ut.Add("phone", "{0} 必须是一个有效的手机号", false)
		case "en":
			e = ut.Add("phone", "{0} must be a valid phone number", false)
		default:
			e = ut.Add("phone", "{0} 必须是一个有效的手机号", false)
		}
		return e
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
}
