package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindMyPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindMyPageLogic(r.Context(), svcCtx)
		resp, err := l.FindMyPage()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindMyPageHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
