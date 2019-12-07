// Package paths is used to give static typing to paths served by the application.
package paths

import (
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	Home                = "home"
	Login               = "login"
	Logout              = "logout"
	Register            = "register"
	ProblemList         = "problem_list"
	Problem             = "problem"
	ProblemNameArg      = "problem_name"
	ProblemFile         = "problem_file"
	ProblemLangArg      = "problem_lang"
	ProblemFileArg      = "problem_file_name"
	SubmitProblem       = "submit_problem"
	Submission          = "submission"
	SubmissionIdArg     = "submission_id"
	User                = "user"
	UserNameArg         = "user_name"
	ContestTeams        = "contest_teams"
	ContestTeamRegister = "contest_team_register"
	ContestTeamStart    = "contest_team_start"
	ContestScoreboard   = "contest_scoreboard"
)

var router = mux.NewRouter()

func GetRouter() *mux.Router {
	return router
}

// Route resolves a route name together with arguments into a URL that can be used to invoke that route.
func Route(name string, args ...interface{}) string {
	var stringified []string
	for _, arg := range args {
		switch a := arg.(type) {
		case string:
			stringified = append(stringified, a)
		case int:
			stringified = append(stringified, strconv.FormatInt(int64(a), 10))
		case int8:
			stringified = append(stringified, strconv.FormatInt(int64(a), 10))
		case int16:
			stringified = append(stringified, strconv.FormatInt(int64(a), 10))
		case int32:
			stringified = append(stringified, strconv.FormatInt(int64(a), 10))
		case int64:
			stringified = append(stringified, strconv.FormatInt(int64(a), 10))
		case uint:
			stringified = append(stringified, strconv.FormatUint(uint64(a), 10))
		case uint8:
			stringified = append(stringified, strconv.FormatUint(uint64(a), 10))
		case uint16:
			stringified = append(stringified, strconv.FormatUint(uint64(a), 10))
		case uint32:
			stringified = append(stringified, strconv.FormatUint(uint64(a), 10))
		case uint64:
			stringified = append(stringified, strconv.FormatUint(uint64(a), 10))
		default:
			panic(fmt.Errorf("used unknown type %T in route", a))
		}
	}
	url, err := GetRouter().Get(name).URL(stringified...)
	if err != nil {
		panic(err)
	}
	return url.String()
}
