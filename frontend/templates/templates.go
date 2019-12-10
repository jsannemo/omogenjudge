package templates

import (
	"fmt"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"html/template"
	"net/http"
	"time"

	"github.com/Masterminds/sprig"

	"github.com/jsannemo/omogenjudge/frontend/util"
)

type TemplateExecutor interface {
	Template() *template.Template
}

var tpls = []string{
	"frontend/templates/*.tpl",
	"frontend/templates/contests/*.tpl",
	"frontend/templates/helpers/*.tpl",
	"frontend/templates/home/*.tpl",
	"frontend/templates/problems/*.tpl",
	"frontend/templates/users/*.tpl",
	"frontend/templates/submissions/*.tpl",
}

func templates() *template.Template {
	tpl := template.New("templates").Funcs(sprig.FuncMap()).Funcs(
		map[string]interface{}{
			"language": util.GetLanguage,
			"durationToSeconds": func(dur time.Duration) string {
				secs := dur.Truncate(time.Second) / time.Second
				return fmt.Sprintf("%d",secs)
			},
			"interval": func(dur time.Duration) string {
				secs := dur.Truncate(time.Second) / time.Second
				return fmt.Sprintf("%02d:%02d:%02d", secs/3600, (secs/60)%60, secs%60)
			},
			"hhmm": func(dur time.Duration) string {
				secs := dur.Truncate(time.Second) / time.Second
				return fmt.Sprintf("%02d:%02d", secs/3600, (secs/60)%60)
			},
			"path": paths.Route,
		})
	for _, t := range tpls {
		tpl = template.Must(tpl.ParseGlob(t))
	}
	return tpl
}

type refreshingExecutor struct{}

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

func tostring(val interface{}) string {
	switch val := val.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

func dict(vals ...interface{}) map[string]interface{} {
	dict := map[string]interface{}{}
	for i := 0; i < len(vals); i += 2 {
		k := tostring(vals[i])
		dict[k] = vals[i+1]
	}
	return dict
}
