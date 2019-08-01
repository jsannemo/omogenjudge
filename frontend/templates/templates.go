package templates

import(
  "html/template"
  "net/http"
)

type TemplateExecutor interface {
  Template() *template.Template
}

func templates() *template.Template {
  return template.Must(template.Must(template.ParseGlob("frontend/templates/*.tpl")).ParseGlob("frontend/templates/**/*.tpl"))
}

type refreshingExecutor struct {}

func (re *refreshingExecutor) Template() *template.Template {
  return templates()
}

type cachingExecutor struct {
  template *template.Template
}

func (ce *cachingExecutor) Template() *template.Template {
  return ce.template
}

// TODO: add some dev env setting to make this caching
var executor = &refreshingExecutor{}

func ExecuteTemplates(w http.ResponseWriter, name string, data interface{}) error {
  tpl := executor.Template()
  if err := tpl.ExecuteTemplate(w, "header", data); err != nil {
    return err
  }
  if err := tpl.ExecuteTemplate(w, "nav", data); err != nil {
    return err
  }
  if err := tpl.ExecuteTemplate(w, name, data); err != nil {
    return err
  }
  if err := tpl.ExecuteTemplate(w, "footer", data); err != nil {
    return err
  }
  return nil
}
