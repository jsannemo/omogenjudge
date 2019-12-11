package problems

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/request"
	"github.com/jsannemo/omogenjudge/frontend/util"
	"github.com/jsannemo/omogenjudge/storage/models"
	"github.com/jsannemo/omogenjudge/storage/problems"
	"github.com/jsannemo/omogenjudge/storage/submissions"
)

type SubmitParams struct {
	Problem   *models.Problem
	Languages []*util.Language
}

func SubmitHandler(r *request.Request) (request.Response, error) {
	team := r.Context.Team
	contest := r.Context.Contest
	if contest != nil && !r.Context.Contest.Started(team) {
		return request.Redirect(paths.Route(paths.Home)), nil
	}

	vars := mux.Vars(r.Request)
	problem, err := getProblemIfVisible(r, vars[paths.ProblemNameArg], problems.ListArgs{WithStatements: problems.StmtTitles})
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return request.NotFound(), nil
	}
	if contest != nil {
		if !contest.HasProblem(problem.ProblemID) {
			return request.NotFound(), nil
		}
	} else if problem.Hidden() {
		return request.NotFound(), nil
	}

	if r.Request.Method == http.MethodPost {
		submit := r.Request.FormValue("submission")
		language := r.Request.FormValue("language")
		l := util.GetLanguage(language)
		if l == nil {
			return request.NotFound(), nil
		}
		s := &models.Submission{
			AccountID: r.Context.UserID,
			ProblemID: problem.ProblemID,
			Language:  l.LanguageId,
			Files: []*models.SubmissionFile{{
				Path:     l.DefaultFile(),
				Contents: submit,
			}},
			CurrentRun: &models.SubmissionRun{
				ProblemVersionID: problem.CurrentVersion.ProblemVersionID,
				Status:           models.StatusNew,
				Evaluation: models.Evaluation{
					Verdict: models.VerdictUnjudged,
				},
			},
		}
		err := submissions.CreateSubmission(r.Request.Context(), s, problem.CurrentVersion.ProblemVersionID)
		if err != nil {
			return request.Error(err), nil
		}
		subUrl := paths.Route(paths.Submission, paths.SubmissionIdArg, s.SubmissionID)
		return request.Redirect(subUrl), nil
	}

	return request.Template("problems_submit", &SubmitParams{problem, util.Languages()}), nil
}
