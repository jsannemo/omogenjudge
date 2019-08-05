// Language-related utilities.
package language

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/logger"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
	"github.com/jsannemo/omogenjudge/runner/compilers"
	"github.com/jsannemo/omogenjudge/runner/runners"
	execpb "github.com/jsannemo/omogenjudge/sandbox/api"
)

var languages = make(map[string]*Language)

// CompileFunc is a function used to compile a program into the given path.
// It may use calls to the execution service in order to perform the compilation.
type CompileFunc func(program *runpb.Program, outputPath string, client execpb.ExecuteServiceClient) (*compilers.Compilation, error)

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
	Program runners.RunFunc
}

// ToApiLanguage converts the internal language to the external API representaiton.
func (l *Language) ToApiLanguage() *runpb.Language {
	return &runpb.Language{
		LanguageId: l.Id,
		Version:    l.Version,
		Group:      l.LanguageGroup,
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

func hasExt(p *runpb.Program, ext string) bool {
	for _, s := range p.Sources {
		if filepath.Ext(s.Path) == ext {
			return true
		}
	}
	return false
}

func hasBang(p *runpb.Program, ext string) bool {
	for _, s := range p.Sources {
		lines := strings.Split(s.Contents, "\n")
		if len(lines) != 0 && strings.HasSuffix(lines[0], ext) {
			return true
		}
	}
	return false
}

// TODO: this is terrible
func GuessLanguage(p *runpb.Program) error {
	if p.LanguageId != "" {
		return nil
	}
	logger.Infof("Program: %v", p.Sources)
	if hasExt(p, ".cpp") || hasExt(p, ".cc") {
		p.LanguageId = "gpp17"
	} else if hasExt(p, ".py") {
		if hasBang(p, "python3") {
			p.LanguageId = "python3"
		} else {
			p.LanguageId = "pypy2"
		}
	}
	if p.LanguageId != "" {
		return nil
	}
	return fmt.Errorf("Could not detect language")
}
