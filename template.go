package testreport

import "text/template"

const defaultTemplateMarkdown = `# Test Report

Total: {{.Tests}} ✔️ Passed: {{.Passed}} ⏩ Skipped: {{.Skipped}} ❌ Failed: {{.Failed}} ⏱️ Duration: {{.Duration}}

{{range .PackageResult}}<details><summary>{{.PackageResult.Icon}} 10/37 {{.Name.Path}}<b>{{.Name.Package}}</b> {{.Duration}}</summary>

{{range .Tests}}<blockquote><details><summary>{{.TestResult.Icon}} {{.Name}} {{.Duration}}</summary><blockquote>

{{range .Output}}` + "`" + `{{.Time.Format "15:04:05.000"}}` + "`" + ` {{EscapeMarkdown .Text}}

{{end}}</blockquote></details></blockquote>{{end}}</details>{{end}}
`

func GetTemplate(pathToTemplate string) (tmp *template.Template, err error) {
	tmp = template.New("template").Funcs(template.FuncMap{
		"EscapeHtml":     EscapeHtml,
		"EscapeMarkdown": EscapeMarkdown,
	})
	if pathToTemplate == "" {
		return template.Must(tmp.Parse(defaultTemplateMarkdown)), nil
	}
	return tmp.ParseFiles(pathToTemplate)
}
