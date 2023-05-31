package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

var cookieStore *sessions.CookieStore

func InitCookieStore(sessionKey string) {
	cookieStore = sessions.NewCookieStore([]byte(sessionKey))
}

func getCookieSessionWithAge(r *http.Request, sessionName, path, appEnv string, age int) (session *sessions.Session, err error) {
	defer func() {
		err = errors.Wrap(err, "getCookieSessionWithAge error")
	}()

	session, err = cookieStore.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	secure := appEnv != "local"
	session.Options = &sessions.Options{
		Path:     path,
		MaxAge:   age,
		HttpOnly: true,
		Secure:   secure,
	}

	return session, nil
}

func getCookieSession(r *http.Request, sessionName, path, appEnv string) (*sessions.Session, error) {
	age := 60 * 60 * 24 * 30 // 30 days
	return getCookieSessionWithAge(r, sessionName, path, appEnv, age)
}
