package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/gorilla/sessions"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	lineLoginSessionState = "state"
)

func LineLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := generateLineLoginState()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineLoginHandler error"))
			return
		}

		session, err := getLineLoginSession(r, svcCtx.Config.AppEnv)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineLoginHandler error"))
			return
		}

		session.Values[lineLoginSessionState] = state
		if err := session.Save(r, w); err != nil {
			httpx.Error(w, errors.Wrap(err, "LineLoginHandler error"))
			return
		}

		url, err := makeLINELoginRedirectURL(state, svcCtx.Config)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineLoginHandler error"))
			return
		}

		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineLoginHandler error"))
		} else {
			http.Redirect(w, r, url, http.StatusFound)
		}
	}
}

// LINEログインに必要なstate値を生成
// https://developers.line.biz/ja/docs/line-login/integrate-line-login/#making-an-authorization-request
func generateLineLoginState() (string, error) {
	randomByte := make([]byte, 32)
	if _, err := rand.Read(randomByte); err != nil {
		return "", err
	}

	state := hex.EncodeToString(randomByte)
	return state, nil
}

func getLineLoginSession(r *http.Request, appEnv string) (*sessions.Session, error) {
	return getCookieSession(r, "line_login", "/line", appEnv)
}

func makeLINELoginRedirectURL(state string, cfg config.Config) (string, error) {
	req, err := http.NewRequest("GET", constant.LineLoginBaseURL, nil)
	if err != nil {
		return "", err
	}

	params := req.URL.Query()
	params.Add("response_type", "code")
	params.Add("bot_prompt", "aggressive")
	params.Add("client_id", fmt.Sprint(cfg.LineLoginChannelID))
	params.Add("redirect_uri", config.LineLoginCallbackURL(cfg.FrontURL))
	params.Add("state", state)
	params.Add("scope", "profile openid")
	req.URL.RawQuery = params.Encode()

	url := req.URL.String()
	return url, nil
}
