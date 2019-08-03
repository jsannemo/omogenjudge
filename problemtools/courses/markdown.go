package courses

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"

	"github.com/jsannemo/omogenjudge/problemtools/util"
)

type markdownResult struct {
	title   string
	summary string
	output  string
}

func handleProblem(lines []string, w io.Writer) error {
	if len(lines) == 0 {
		return fmt.Errorf("missing problem ID")
	}
	tpl := fmt.Sprintf(`{{ template "course_helper_problem" "%s" }}`, lines[0])
	_, err := w.Write([]byte(tpl))
	return err
}

func handleSummary(lines []string, res *markdownResult) error {
	if len(lines) == 0 {
		return fmt.Errorf("missing summary")
	}
	res.summary = lines[0]
	return nil
}

func handleBoxStart(heading string, w io.Writer) error {
	tpl := fmt.Sprintf(`
  <div class="course-box">
    <div class="course-box-header">%s</div>
    <div class="course-box-content">`, heading)
	_, err := w.Write([]byte(tpl))
	return err
}

func handleBoxEnd(w io.Writer) error {
	tpl := "</div></div>"
	_, err := w.Write([]byte(tpl))
	return err
}

func processInlineCommands(cmd string, w io.Writer) error {
	cmd = strings.TrimPrefix(cmd, "omogen ")
	if strings.HasPrefix(cmd, "box ") {
		cmd = strings.TrimPrefix(cmd, "box ")
		if err := handleBoxStart(cmd, w); err != nil {
			return err
		}
	} else if cmd == "box" {
		if err := handleBoxEnd(w); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown inline command: %v", cmd)
	}
	return nil
}

func processCommands(res *markdownResult, cmd string, w io.Writer) error {
	lines := strings.Split(strings.TrimSpace(cmd), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	switch lines[0] {
	case "problem":
		if err := handleProblem(lines[1:], w); err != nil {
			return err
		}
	case "summary":
		if err := handleSummary(lines[1:], res); err != nil {
			return err
		}
	}
	return nil
}

func parseMarkdown(path string, rep util.Reporter) (*markdownResult, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	inTitle := false
	res := &markdownResult{}

	opts := html.RendererOptions{
		Flags: html.CommonFlags,
		RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			if code, ok := node.(*ast.CodeBlock); ok {
				if code.IsFenced && string(code.Info) == "omogen" {
					err = processCommands(res, string(code.Literal), w)
					if err != nil {
						return ast.Terminate, true
					}
					return ast.GoToNext, true
				}
			}
			if code, ok := node.(*ast.Code); ok && strings.HasPrefix(string(code.Literal), "omogen ") {
				err = processInlineCommands(string(code.Literal), w)
				if err != nil {
					return ast.Terminate, true
				}
				return ast.GoToNext, true
			}

			if heading, ok := node.(*ast.Heading); ok {
				if entering {
					inTitle = res.title == "" && heading.Level == 1
				} else if inTitle {
					inTitle = false
					return ast.GoToNext, true
				}
			}
			if inTitle {
				if leaf := node.AsLeaf(); leaf != nil {
					res.title = res.title + string(leaf.Literal)
				}
				return ast.GoToNext, true
			}
			return ast.GoToNext, false
		},
	}

	renderer := html.NewRenderer(opts)
	if err != nil {
		rep.Err("failed parsing section: %v", err)
		return nil, nil
	}
	res.output = string(markdown.ToHTML(dat, nil, renderer))
	return res, nil
}
