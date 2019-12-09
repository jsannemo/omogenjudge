package main

import (
	"fmt"

	"github.com/jsannemo/omogenjudge/frontend/contests"
	"github.com/jsannemo/omogenjudge/frontend/home"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/problems"
	"github.com/jsannemo/omogenjudge/frontend/submissions"
	"github.com/jsannemo/omogenjudge/frontend/users"

	"github.com/gorilla/mux"
)

func configureRouter() *mux.Router {
	r := paths.GetRouter()
	registerStaticHandler(r)
	r.HandleFunc("/", plain(home.HomeHandler)).Name(paths.Home)
	r.HandleFunc("/register", plain(users.RegisterHandler)).Name(paths.Register)
	r.HandleFunc("/login", plain(users.LoginHandler)).Name(paths.Login)
	r.HandleFunc("/logout", plain(users.LogoutHandler)).Name(paths.Logout)
	r.HandleFunc("/problems", plain(problems.ListHandler)).Name(paths.ProblemList)
	r.HandleFunc(fmt.Sprintf("/problems/{%s}", paths.ProblemNameArg), plain(problems.ViewHandler)).Name(paths.Problem)
	r.HandleFunc(fmt.Sprintf("/problems/{%s}/submit", paths.ProblemNameArg), plain(problems.SubmitHandler)).Name(paths.SubmitProblem)
	r.HandleFunc(fmt.Sprintf("/problems/{%s}/{%s}/{%s}",
		paths.ProblemNameArg, paths.ProblemLangArg, paths.ProblemFileArg),
		plain(problems.FileHandler)).Name(paths.ProblemFile)
	r.HandleFunc(fmt.Sprintf("/submissions/{%s}", paths.SubmissionIdArg), plain(submissions.ViewHandler)).Name(paths.Submission)
	r.HandleFunc(fmt.Sprintf("/users/{%s}", paths.UserNameArg), plain(users.ViewHandler)).Name(paths.User)
	r.HandleFunc("/teams", plain(contests.TeamListHandler)).Name(paths.ContestTeams)
	r.HandleFunc("/teams/register", plain(contests.RegisterHandler)).Name(paths.ContestTeamRegister)
	r.HandleFunc("/teams/start", plain(contests.StartHandler)).Name(paths.ContestTeamStart)
	r.HandleFunc("/scoreboard", plain(contests.ScoreboardHandler)).Name(paths.ContestScoreboard)
	return r
}
