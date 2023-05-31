package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindAdminRoleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindAdminRoleListLogic(r.Context(), svcCtx)
		resp, err := l.FindAdminRoleList(r.Context())
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindAdminRoleListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
