// Sets up routing of all request paths to their handlers.
package main

import (
	"fmt"

	"github.com/jsannemo/omogenjudge/frontend/courses"
	"github.com/jsannemo/omogenjudge/frontend/paths"
	"github.com/jsannemo/omogenjudge/frontend/problems"
	"github.com/jsannemo/omogenjudge/frontend/submissions"
	"github.com/jsannemo/omogenjudge/frontend/users"

	"github.com/gorilla/mux"
)

func configureRouter() *mux.Router {
	r := paths.GetRouter()
	registerStaticHandler(r)
	// TODO update when there is a home handler
	r.HandleFunc("/", plain(problems.ListHandler)).Name(paths.Home)
	r.HandleFunc("/register", plain(users.RegisterHandler)).Name(paths.Register)
	r.HandleFunc("/login", plain(users.LoginHandler)).Name(paths.Login)
	r.HandleFunc("/logout", plain(users.LogoutHandler)).Name(paths.Logout)
	r.HandleFunc("/problems", plain(problems.ListHandler)).Name(paths.ProblemList)
	r.HandleFunc(fmt.Sprintf("/problems/{%s}", paths.ProblemNameArg), plain(problems.ViewHandler)).Name(paths.Problem)
	r.HandleFunc(fmt.Sprintf("/problems/{%s}/submit", paths.ProblemNameArg), plain(problems.SubmitHandler)).Name(paths.SubmitProblem)
	r.HandleFunc(fmt.Sprintf("/submissions/{%s}", paths.SubmissionIdArg), plain(submissions.ViewHandler)).Name(paths.Submission)
	r.HandleFunc(fmt.Sprintf("/users/{%s}", paths.UserNameArg), plain(users.ViewHandler)).Name(paths.User)
	r.HandleFunc("/courses", plain(courses.ListHandler)).Name(paths.CourseList)
	r.HandleFunc(fmt.Sprintf("/courses/{%s}", paths.CourseNameArg), plain(courses.CourseHandler)).Name(paths.Course)
	r.HandleFunc(fmt.Sprintf("/courses/{%s}/{%s}", paths.CourseNameArg, paths.CourseChapterNameArg), plain(courses.ChapterHandler)).Name(paths.CourseChapter)
	r.HandleFunc(fmt.Sprintf("/courses/{%s}/{%s}/{%s}", paths.CourseNameArg, paths.CourseChapterNameArg, paths.CourseSectionNameArg), plain(courses.SectionHandler)).Name(paths.CourseSection)
	return r
}
