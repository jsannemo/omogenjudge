package language

import (
	"fmt"
	"path/filepath"
	"strings"

	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

// TODO: this should be made more robust and a bit more abstract.

// GuessLanguage tries to detect the language of the given uncompiled program.
func GuessLanguage(p *runpb.Program) error {
	if p.LanguageId != "" {
		return nil
	}
	if hasExt(p, ".cpp") || hasExt(p, ".cc") {
		p.LanguageId = "gpp17"
	} else if hasExt(p, ".py") {
		if hasBang(p, "python3") {
			p.LanguageId = "cpython3"
		} else {
			p.LanguageId = "pypy2"
		}
	}
	if p.LanguageId != "" {
		return nil
	}
	return fmt.Errorf("Could not detect language")
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
