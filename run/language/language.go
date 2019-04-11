package language

import (
  "log"

	execpb "github.com/jsannemo/omogenexec/exec/api"
	runpb "github.com/jsannemo/omogenexec/run/api"
)

var languages = make(map[string]*Language)

type CompileFunc func(*runpb.Program, string, execpb.ExecuteServiceClient) (*runpb.CompiledProgram, error)
type RunFunc func(*runpb.RunRequest, execpb.ExecuteService_ExecuteClient) (*runpb.RunResponse, error)

type Language struct {
  Id string
	Version string
  LanguageGroup runpb.LanguageGroup
  Compile CompileFunc
  Run RunFunc
}

func (l *Language) ToApiLanguage() *runpb.Language {
  return &runpb.Language{
    LanguageId: l.Id,
    Version: l.Version,
    Group: l.LanguageGroup,
  }
}

func GetLanguages() map[string]*Language {
  return languages
}

func GetLanguage(id string) (*Language, bool) {
  lang, exists := languages[id]
  return lang, exists
}

func registerLanguage(language *Language) {
  log.Printf("Registering language: %v", *language)
  languages[language.Id] = language
}
