// paths is used to give static typing to paths served by the application.
package paths

import (
	"github.com/gorilla/mux"
)

const (
	Home           = "home"
	Login          = "login"
	Logout         = "logout"
	Register       = "register"
	ProblemList    = "problem_list"
	Problem        = "problem"
	ProblemNameArg = "problem_name"
  SubmitProblem  = "submit_problem"
)

var router = mux.NewRouter()

func GetRouter() *mux.Router {
	return router
}

func Route(name string) *mux.Route {
	return GetRouter().Get(name)
}
