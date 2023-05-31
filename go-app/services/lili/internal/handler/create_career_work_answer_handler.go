package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateCareerWorkAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCareerWorkAnswerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateCareerWorkAnswerHandler error"))
			return
		}

		l := logic.NewCreateCareerWorkAnswerLogic(r.Context(), svcCtx)
		err := l.CreateCareerWorkAnswer(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateCareerWorkAnswerHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
