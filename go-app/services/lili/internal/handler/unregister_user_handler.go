package handler

import (
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnregisterUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnregisterUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UnregisterUserHandler error"))
			return
		}

		l := logic.NewUnregisterUserLogic(r.Context(), svcCtx)
		err := l.UnregisterUser(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UnregisterUserHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
