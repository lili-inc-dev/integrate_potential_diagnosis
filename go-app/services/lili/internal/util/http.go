package util

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func SetReqContentTypeURLEncoded(req *http.Request) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func SetReqContentTypeJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func SetResContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func GetBearerToken(req *http.Request) (token string, err error) {
	defer func() {
		err = errors.Wrap(err, "GetBearerToken error")
	}()

	authHeader := req.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return "", errors.New("authorization header format is invalid")
	}

	return splitToken[1], nil
}

func SetBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}
