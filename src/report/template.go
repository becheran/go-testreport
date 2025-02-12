package report

import (
	_ "embed"
	"path/filepath"
	"text/template"
)

//go:embed templates/md.tmpl
var defaultTemplateMarkdown string

func GetTemplate(pathToTemplate string) (tmp *template.Template, err error) {
	tmp = template.New(filepath.Base(pathToTemplate)).Funcs(template.FuncMap{
		"EscapeHtml":     EscapeHtml,
		"EscapeMarkdown": EscapeMarkdown,
	})
	if pathToTemplate == "" {
		return template.Must(tmp.Parse(defaultTemplateMarkdown)), nil
	}
	return tmp.ParseFiles(pathToTemplate)
}
