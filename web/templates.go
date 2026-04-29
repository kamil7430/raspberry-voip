package web

import (
	"embed"
	"html/template"
	"log"
)

//go:embed templates/*.html
var assets embed.FS

// Templates is a pointer to template.Template struct of loaded HTML templates from /web/templates directory.
var Templates = func() *template.Template {
	templates, err := template.ParseFS(assets, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	return templates
}()
