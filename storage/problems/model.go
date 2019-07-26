// Representation of the problem storage model
package problems

import (
  "html/template"

  "github.com/google/logger"
  "golang.org/x/text/language"

  "github.com/jsannemo/omogenjudge/frontend/paths"
  "github.com/jsannemo/omogenjudge/storage/files"
)

type ProblemMap map[int32]*Problem
type TestGroupMap map[int32]*TestCaseGroup

type ProblemStatement struct {
  // The language tag for the statement.
  Language string

  // The title of the statement.
  Title string

  // The HTML template of the statement.
  Html template.HTML
}

type Problem struct {
  // A numeric problem ID.
  // This should not be exposed externally.
  ProblemId int32

  // The short name of the problem.
  // This is suitable to use in e.g. URLs or as externally-visible identifiers.
  ShortName string

  // A list of all statements corresponding to a problem.
  Statements []*ProblemStatement

  TestGroups []*TestCaseGroup
}

type TestCaseGroup struct {
  // A numeric test case group ID.
  // This should not be exposed externally.
  TestCaseGroupId int32

  // The name of the test group.
  // This can be exposed externally.
  Name string

  // Whether the test cases of this group are publically visisible.
  PublicVisibility bool

  // The test cases this group contains.
  Tests []*TestCase
}

type TestCase struct {
  // A numeric test case ID.
  // This should not be exposed externally.
  TestCaseId int32

  // The name of the test case.
  // This should not be exposed externally.
  Name string

  InputFile *files.StoredFile

  OutputFile *files.StoredFile
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

// LocalizedTitle returns the title of the statement best corresponding to the languages in langs.
func (p *Problem) LocalizedTitle(langs []language.Tag) string {
  return localizedStatement(p, langs).Title
}

// LocalizedStatement returns the HTML statement of the statement best corresponding to the languages in langs.
func (p *Problem) LocalizedStatement(langs []language.Tag) template.HTML {
  return localizedStatement(p, langs).Html
}

// Link returns the link to this problem.
func (p *Problem) Link() string {
  link, err := paths.Route(paths.Problem).URL(paths.ProblemNameArg, p.ShortName)
  if err != nil {
    logger.Fatalf("Faild creating problem link: %v", err)
  }
  return link.String()
}

// Link returns the link to this problem.
func (p *Problem) SubmitLink() string {
  link, err := paths.Route(paths.SubmitProblem).URL(paths.ProblemNameArg, p.ShortName)
  if err != nil {
    logger.Fatalf("Faild creating submit problem link: %v", err)
  }
  return link.String()
}

// Samples returns all testcases that should be publically visible
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

func (p *Problem) TestDataFiles() files.FileList {
  var files files.FileList
  for _, group := range p.TestGroups {
    for _, tc := range group.Tests {
      files = append(files, tc.InputFile, tc.OutputFile)
    }
  }
  return files
}

func (p *Problem) Tests() []*TestCase {
  var tests []*TestCase
  for _, group := range p.TestGroups {
    for _, tc := range group.Tests {
      tests = append(tests, tc)
    }
  }
  return tests
}
