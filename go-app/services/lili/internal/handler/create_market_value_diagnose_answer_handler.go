package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateMarketValueDiagnoseAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateMarketValueDiagnoseAnswerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateMarketValueDiagnoseAnswerHandler error"))
			return
		}

		l := logic.NewCreateMarketValueDiagnoseAnswerLogic(r.Context(), svcCtx)
		err := l.CreateMarketValueDiagnoseAnswer(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateMarketValueDiagnoseAnswerHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
