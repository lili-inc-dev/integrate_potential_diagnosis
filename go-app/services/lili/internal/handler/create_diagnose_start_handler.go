package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateDiagnoseStartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCreateDiagnoseStartLogic(r.Context(), svcCtx)
		err := l.CreateDiagnoseStart()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateDiagnoseStartHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
