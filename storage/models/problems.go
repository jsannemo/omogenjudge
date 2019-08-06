package models

import (
	"database/sql"
	"fmt"
	"html/template"

	"golang.org/x/text/language"

	"github.com/jsannemo/omogenjudge/frontend/paths"
)

type ProblemMap map[int32]*Problem

func (p ProblemMap) AsList() ProblemList {
	var probs ProblemList
	for _, prob := range p {
		probs = append(probs, prob)
	}
	return probs
}

func (p ProblemMap) Ids() []int32 {
	var ids []int32
	for id, _ := range p {
		ids = append(ids, id)
	}
	return ids
}

type ProblemList []*Problem

func (pl ProblemList) AsMap() ProblemMap {
	pm := make(ProblemMap)
	for _, p := range pl {
		pm[p.ProblemId] = p
	}
	return pm
}

type StatementList []*ProblemStatement

func (sl StatementList) AddTo(pm ProblemMap) {
	for _, s := range sl {
		p := pm[s.ProblemId]
		p.Statements = append(p.Statements, s)
	}
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

type Problem struct {
	// A numeric problem ID.
	// This should not be exposed externally.
	ProblemId int32 `db:"problem_id"`

	// The short name of the problem.
	// This is suitable to use in e.g. URLs or as externally-visible identifiers.
	ShortName string `db:"short_name"`

	// A list of all statements corresponding to a problem.
	Statements []*ProblemStatement

	TestGroups []*TestGroup

	Author string

	License License

	TimeLimMs int32 `db:"time_limit_ms"`

	MemLimKb int32 `db:"memory_limit_kb"`

	OutputValidator *OutputValidator `db:"problem_output_validator"`
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

func (p *Problem) Samples() []*TestCase {
	var samples []*TestCase
	for _, group := range p.TestGroups {
		if !group.PublicVisibility {
			continue
		}
		samples = append(samples, group.Tests...)
	}
	return samples
}

func (p *Problem) Tests() []*TestCase {
	var tests []*TestCase
	for _, group := range p.TestGroups {
		tests = append(tests, group.Tests...)
	}
	return tests
}

func (p *Problem) TestDataFiles() FileList {
	var files FileList
	for _, tc := range p.Tests() {
		files = append(files, tc.InputFile, tc.OutputFile)
	}
	return files
}

func (p *Problem) TimeLimString() string {
	return fmt.Sprintf("%.1g s", float64(p.TimeLimMs)/1000)
}

func (p *Problem) MemLimString() string {
	return fmt.Sprintf("%.1g GB", float64(p.MemLimKb)/1000/1000)
}

type OutputValidator struct {
	ValidatorLanguageId sql.NullString      `db:"language_id"`
	ValidatorSourceZip  *NullableStoredFile `db:"validator_source_zip"`
}

func (ov *OutputValidator) Nil() bool {
	return ov.ValidatorSourceZip.NotNil()
}
