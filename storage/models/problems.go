package models

import (
	"html/template"

	"golang.org/x/text/language"

	"github.com/jsannemo/omogenjudge/frontend/paths"
)

type Problem struct {
	// A numeric problem ID.
	// This should not be exposed externally.
	ProblemId int32 `db:"problem_id"`

	// The short name of the problem.
	// This is suitable to use in e.g. URLs or as externally-visible identifiers.
	ShortName string `db:"short_name"`

	// A list of all statements corresponding to a problem.
	Statements []*ProblemStatement

	Author string `db:"author"`

	License License `db:"license"`

	CurrentVersion *ProblemVersion `db:"problem_version"`
}

// localizedStatement returns the statement of a problem closest to the ones given in langs.
// By default, "en" and "sv" are fallback languages
func localizedStatement(p *Problem, langs []language.Tag) *ProblemStatement {
	var has []language.Tag
	userPrefs := append(langs, language.Make("en"), language.Make("sv"))
	for _, statement := range p.Statements {
		has = append(has, language.Make(statement.Language))
	}
	var matcher = language.NewMatcher(has)
	_, index, _ := matcher.Match(userPrefs...)
	return p.Statements[index]
}

func (p *Problem) LocalizedTitle(preferred []language.Tag) string {
	return localizedStatement(p, preferred).Title
}

func (p *Problem) LocalizedStatement(preferred []language.Tag) template.HTML {
	return template.HTML(localizedStatement(p, preferred).Html)
}

func (p *Problem) Link() string {
	return paths.Route(paths.Problem, paths.ProblemNameArg, p.ShortName)
}

func (p *Problem) SubmitLink() string {
	return paths.Route(paths.SubmitProblem, paths.ProblemNameArg, p.ShortName)
}

type ProblemStatement struct {
	ProblemId int32 `db:"problem_id"`
	// The language tag for the statement.
	Language string `db:"language"`

	// The title of the statement.
	Title string `db:"title"`

	// The HTML template of the statement.
	Html string `db:"html"`
}
