// Language-related utilities.
package language

import (
  "github.com/google/logger"

	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

var languages = make(map[string]*Language)

// CompileFunc is a function used to compile a program into the given path.
// It may use calls to the execution service in order to perform the compilation.
type CompileFunc func(program *runpb.Program, outputPath string, client execpb.ExecuteServiceClient) (*runpb.CompiledProgram, error)

// RunFunc is a function used to run a given program.
type RunFunc func(*runpb.RunRequest, execpb.ExecuteService_ExecuteClient) (*runpb.RunResponse, error)

// A programming language.
type Language struct {
  // An identifier for this language.
  // This is suitable for inclusion in URLs, and can be displayed externally.
  Id string

  // The version that this language belongs to.
	Version string

  // The language group that this language belongs to.
  LanguageGroup runpb.LanguageGroup

  // The compile function that should be used to compile programs of this language.
  Compile CompileFunc

  // The run function that should be used to run compiled programs of this language.
  Run func() RunFunc
}

// ToApiLanguage converts the internal language to the external API representaiton.
func (l *Language) ToApiLanguage() *runpb.Language {
  return &runpb.Language{
    LanguageId: l.Id,
    Version: l.Version,
    Group: l.LanguageGroup,
  }
}

// GetLanguages returns all installed languages, mapped from language ID to the language itself.
func GetLanguages() map[string]*Language {
  return languages
}

// GetLanguage returns the language with the given id.
// The second parameter designates whether a language with this ID existed.
func GetLanguage(id string) (lang *Language, found bool) {
  lang, found = languages[id]
  return
}

func registerLanguage(language *Language) {
  logger.Infof("Registering language: %v", *language)
  languages[language.Id] = language
}
