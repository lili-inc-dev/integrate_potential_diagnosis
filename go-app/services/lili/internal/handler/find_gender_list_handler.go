package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindGenderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindGenderListLogic(r.Context(), svcCtx)
		resp, err := l.FindGenderList()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindGenderListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
