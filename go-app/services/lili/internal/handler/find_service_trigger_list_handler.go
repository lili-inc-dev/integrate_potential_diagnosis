package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindServiceTriggerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindServiceTriggerListLogic(r.Context(), svcCtx)
		resp, err := l.FindServiceTriggerList()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindServiceTriggerListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
