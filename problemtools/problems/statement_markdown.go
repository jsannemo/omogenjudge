package problems

import (
	"fmt"
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

func parseMarkdown(path string, problemName string, reporter util.Reporter) (*toolspb.ProblemStatement, error) {
	dat, err := ioutil.ReadFile(filepath.Join(path, "problem.md"))
	if err != nil {
		return nil, err
	}

	title := ""
	inTitle := false
	lang := filepath.Base(path)

	// Extract the first title tag (and remove it)
	opts := html.RendererOptions{
		Flags: html.CommonFlags,
		RenderNodeHook: func(_ io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			if img, ok := node.(*ast.Image); ok {
				img.Destination = []byte(fmt.Sprintf("/problems/%s/%s/%s", problemName, lang, string(img.Destination)))
			}
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
	statement_html := string(markdown.ToHTML(dat, nil, renderer))

	if title == "" {
		reporter.Err("Statement for language %s had no title", lang)
	}

	otherFiles := make(map[string]string)
	subFiles, err := listSubFiles(path)
	if err != nil {
		return nil, err
	}
	for _, filePath := range subFiles {
		name := filepath.Base(filePath)
		if name != "problem.md" && filePath != path {
			otherFiles[filepath.Base(filePath)] = filePath
		}
	}

	return &toolspb.ProblemStatement{
		LanguageCode:   lang,
		Title:          title,
		StatementHtml:  statement_html,
		StatementFiles: otherFiles,
	}, nil

}
