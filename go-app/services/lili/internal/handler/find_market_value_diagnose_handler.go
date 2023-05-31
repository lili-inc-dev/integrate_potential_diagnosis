package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindMarketValueDiagnoseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindMarketValueDiagnoseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindMarketValueDiagnoseHandler error"))
			return
		}

		l := logic.NewFindMarketValueDiagnoseLogic(r.Context(), svcCtx)
		resp, err := l.FindMarketValueDiagnose(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindMarketValueDiagnoseHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
