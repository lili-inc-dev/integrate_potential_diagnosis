package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindNoticeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindNoticeListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeListHandler error"))
			return
		}

		l := logic.NewFindNoticeListLogic(r.Context(), svcCtx)
		resp, err := l.FindNoticeList(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
