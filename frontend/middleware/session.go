// Session-related middlewares.
package middleware

import (
	"flag"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/jsannemo/omogenjudge/frontend/request"
)

var (
	cookieSecret = flag.String("cookie_secret", "TMP", "A string used to encrypt cookies - should be random and long enough to be unbruteforcable")
)

var store *sessions.CookieStore

func cookieStore() *sessions.CookieStore {
	if store == nil {
		store = sessions.NewCookieStore([]byte(*cookieSecret))
	}
	return store
}

const (
	sessionKey = "omogenjudge_session"
	userIdKey  = "omogenjudge_session_userid"
)

// readSession reads the session data from the request cookie
// Context values stored in cookies are also inserted in the context.
func readSession(r *request.Request) (request.Response, error) {
	session, err := cookieStore().Get(r.Request, sessionKey)
	if err != nil {
		return nil, err
	}
	r.Session = session

	userId, had := session.Values[userIdKey]
	if had {
		id, err := strconv.Atoi(userId.(string))
		if err != nil {
			return nil, err
		}
		r.Context.UserId = int32(id)
	}
	return nil, nil
}

// writeSession writes the session data to the request cookies.
// Context values that should be persisted are also stored in the cookies.
func writeSession(r *request.Request) (request.Response, error) {
	userId := r.Context.UserId
	if userId != 0 {
		r.Session.Values[userIdKey] = strconv.Itoa(int(userId))
	} else {
		delete(r.Session.Values, userIdKey)
	}
	if err := r.Session.Save(r.Request, r.Writer); err != nil {
		return nil, err
	}
	return nil, nil
}
