package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LineCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LineCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "LineCallbackHandler error"))
			return
		}

		session, err := getLineLoginSession(r, svcCtx.Config.AppEnv)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineCallbackHandler error"))
			return
		}

		stateExpected, ok := session.Values[lineLoginSessionState].(string)
		if !ok {
			err := errors.New("type missmatch of state")
			httpx.Error(w, errors.Wrap(err, "LineCallbackHandler error"))
			return
		}

		if req.State != stateExpected {
			err := errors.New("invalid state value")
			httpx.Error(w, errors.Wrap(err, "LineCallbackHandler error"))
			return
		}

		l := logic.NewLineCallbackLogic(r.Context(), svcCtx)
		resp, err := l.LineCallback(req.AuthCode)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "LineCallbackHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
