// Middlewares are used to execute code before and after a request handler.
package middleware

import (
	"github.com/jsannemo/omogenjudge/frontend/request"
)

// A processor takes a request and returns an optional Response.
// The processor may modify the request values.
type Processor func(*request.Request) (request.Response, error)

// A middleware is a processor with some associated metadata.
type Middleware struct {
	// The processor run as part of the middlware
	Processor Processor

	// Normally, middleware is not run after a prior middleware has generated a response.
	// This flag can be used to override this behaviour and always run a middleware.
	Always bool
}

// WithMiddleware takes a Processor and adds pre- and postprocessing middlware to it.
func WithMiddlewares(responseFn Processor) []Middleware {
	return []Middleware{
		Middleware{readSession, true},
		Middleware{i18n, true},
		Middleware{readUser, false},
		Middleware{responseFn, false},
		Middleware{writeSession, true},
	}
}
