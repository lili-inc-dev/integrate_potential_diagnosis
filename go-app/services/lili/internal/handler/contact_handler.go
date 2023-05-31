package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ContactHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ContactReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "ContactHandler error"))
			return
		}

		l := logic.NewContactLogic(r.Context(), svcCtx)
		err := l.Contact(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "ContactHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
