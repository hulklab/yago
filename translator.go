package yago

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

var uni *ut.UniversalTranslator

func init() {
	cn := zh.New()
	eng := en.New()

	uni = ut.New(cn, cn, eng)
}

func GetLang() string {
	var l string
	if Config.IsSet("app.lang") {
		l = Config.GetString("app.lang")
	} else {
		l = "zh"
	}
	return l
}

func GetTranslator(lang ...string) ut.Translator {
	var l string
	if len(lang) == 0 {
		// 从配置中取
		l = GetLang()
	} else {
		l = lang[0]
	}

	trans, _ := uni.GetTranslator(l)
	return trans
}
