package util

import (
	"context"

	"github.com/google/logger"

	masterpb "github.com/jsannemo/omogenjudge/masterjudge/api"
	masterclient "github.com/jsannemo/omogenjudge/masterjudge/client"
	runpb "github.com/jsannemo/omogenjudge/runner/api"
)

type Language runpb.Language

var langs []*Language

func Languages() []*Language {
	if langs == nil {
		masterc := masterclient.NewClient()
		resp, err := masterc.GetLanguages(context.Background(), &masterpb.GetLanguagesRequest{})
		if err != nil {
			logger.Fatal(err)
		}
		for _, lang := range resp.InstalledLanguages {
			l := Language(*lang)
			langs = append(langs, &l)
		}
	}
	return langs
}

func GetLanguage(tag string) *Language {
	for _, l := range Languages() {
		if l.LanguageId == tag {
			return l
		}
	}
	return nil
}

func (l Language) Name() string {
	switch l.Group {
	case runpb.LanguageGroup_CPP_11:
		return "C++ 11"
	case runpb.LanguageGroup_CPP_14:
		return "C++ 14"
	case runpb.LanguageGroup_CPP_17:
		return "C++ 17"
	case runpb.LanguageGroup_PYTHON_2:
		return "Python 2"
	case runpb.LanguageGroup_PYTHON_2_PYPY:
		return "Python 2 - PyPy"
	case runpb.LanguageGroup_PYTHON_3:
		return "Python 3"
	case runpb.LanguageGroup_PYTHON_3_PYPY:
		return "Python 3 - PyPy"
	}
	return "Unknown language"
}

func (l Language) DefaultFile() string {
	switch l.Group {
	case runpb.LanguageGroup_CPP_11:
		fallthrough
	case runpb.LanguageGroup_CPP_14:
		fallthrough
	case runpb.LanguageGroup_CPP_17:
		return "main.cpp"
	case runpb.LanguageGroup_PYTHON_2:
		fallthrough
	case runpb.LanguageGroup_PYTHON_2_PYPY:
		fallthrough
	case runpb.LanguageGroup_PYTHON_3:
		fallthrough
	case runpb.LanguageGroup_PYTHON_3_PYPY:
		return "main.py"
	default:
		return "main"
	}
}

func (l Language) VsName() string {
	switch l.Group {
	case runpb.LanguageGroup_CPP_11:
		fallthrough
	case runpb.LanguageGroup_CPP_14:
		fallthrough
	case runpb.LanguageGroup_CPP_17:
		return "cpp"
	case runpb.LanguageGroup_PYTHON_2:
		fallthrough
	case runpb.LanguageGroup_PYTHON_2_PYPY:
		fallthrough
	case runpb.LanguageGroup_PYTHON_3:
		fallthrough
	case runpb.LanguageGroup_PYTHON_3_PYPY:
		return "python"
	default:
		return ""
	}
}
