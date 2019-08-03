package problems

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"

	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

func markdownParser() *statementParser {
	return &statementParser{
		"Markdown",
		hasProblemMd,
		parseMarkdown,
	}
}

func listSubFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func hasProblemMd(path string) (bool, error) {
	want := filepath.Join(path, "problem.md")
	files, err := listSubFiles(path)
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if file == want {
			return true, nil
		}
	}
	return false, nil
}

func parseMarkdown(path string, reporter util.Reporter) (*toolspb.ProblemStatement, error) {
	dat, err := ioutil.ReadFile(filepath.Join(path, "problem.md"))
	if err != nil {
		return nil, err
	}

	title := ""
	inTitle := false

	// Extract the first title tag (and remove it)
	opts := html.RendererOptions{
		Flags: html.CommonFlags,
		RenderNodeHook: func(_ io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			ignore := inTitle
			if heading, ok := node.(*ast.Heading); ok {
				if title == "" && entering && heading.Level == 1 {
					inTitle = true
					ignore = true
				}
				if !entering && heading.Level == 1 {
					inTitle = false
				}
			}
			leaf := node.AsLeaf()
			if inTitle && leaf != nil {
				title = title + string(leaf.Literal)
			}
			return ast.GoToNext, ignore
		},
	}

	renderer := html.NewRenderer(opts)
	html := string(markdown.ToHTML(dat, nil, renderer))

	lang := filepath.Base(path)

	if title == "" {
		reporter.Err("Statement for language %s had no title", lang)
	}

	return &toolspb.ProblemStatement{
		LanguageCode:  lang,
		Title:         title,
		StatementHtml: html,
	}, nil

}
