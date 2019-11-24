// Handles processing of an incoming request.
package main

import (
	"errors"
	"net/http"

	"github.com/jsannemo/omogenjudge/frontend/middleware"
	"github.com/jsannemo/omogenjudge/frontend/request"

	"github.com/google/logger"
)

// plain wraps a request handler in the general middleware.
func plain(fn middleware.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := request.NewRequest(w, r)
		middlewares := middleware.WithMiddlewares(fn)
		executeMiddlewares(req, middlewares)
	}
}

// executeMiddlewares executes a request on a list of middlewares, writing the response back to the client
func executeMiddlewares(req *request.Request, middlewares []middleware.Middleware) {
	for _, middleware := range middlewares {
		// Skip further middlewares by default if we have a response
		if !middleware.Always && req.Response != nil {
			continue
		}
		resp, err := middleware.Processor(req)
		if err != nil {
			logger.Errorf("Error during request handling: %v", err)
			// TODO(jsannemo): display actual errors for admins
			req.Response = request.Error(err)
			break
		}
		if resp != nil {
			req.Response = resp
		}
	}
	if req.Response == nil {
		req.Response = request.Error(errors.New("no response generated"))
	}
	req.Write(req.Writer)
}
