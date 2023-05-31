package handler

import (
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FrontUserDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.UserBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserDetailHandler error"))
			return
		}

		var req types.FrontUserDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserDetailHandler error"))
			return
		}

		l := logic.NewFrontUserDetailLogic(r.Context(), svcCtx)
		resp, err := l.FrontUserDetail(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserDetailHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
