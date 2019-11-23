package middleware

import (
	"golang.org/x/text/language"

	"github.com/jsannemo/omogenjudge/frontend/request"
)

// i18n is a middleware that parses the accept languages of the client and stores them in the context.
func i18n(r *request.Request) (request.Response, error) {
	r.Context.Locales, _, _ = language.ParseAcceptLanguage(r.Request.Header.Get("Accept-Language"))
	return nil, nil
}
