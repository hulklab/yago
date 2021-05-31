// +build go1.16

package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/hulklab/yago"
)

//go:embed dist/* dist/assets/*
var f embed.FS

func init() {
	fmt.Println("here")
	yago.AddAppInitHook(func(app *yago.App) error {
		httpEngine := app.HttpEngine()

		tpl := template.Must(template.New("").ParseFS(f, "dist/*.html"))
		httpEngine.SetHTMLTemplate(tpl)

		httpEngine.StaticFS("/public/", http.FS(f))

		return nil
	})
}
