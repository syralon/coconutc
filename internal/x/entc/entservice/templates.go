package entservice

import (
	"embed"
	"text/template"
)

//go:embed templates
var fs embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(fs, "templates/*.tpl")
	if err != nil {
		panic(err)
	}
}
