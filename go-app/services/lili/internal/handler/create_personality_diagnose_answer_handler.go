package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreatePersonalityDiagnoseAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePersonalityDiagnoseAnswerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreatePersonalityDiagnoseAnswerHandler error"))
			return
		}

		l := logic.NewCreatePersonalityDiagnoseAnswerLogic(r.Context(), svcCtx)
		err := l.CreatePersonalityDiagnoseAnswer(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreatePersonalityDiagnoseAnswerHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
