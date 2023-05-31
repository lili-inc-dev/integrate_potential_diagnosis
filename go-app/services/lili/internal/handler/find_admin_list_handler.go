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

func FindAdminListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.AdminBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindAdminListHandler error"))
			return
		}

		var req types.FindAdminListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindAdminListHandler error"))
			return
		}

		l := logic.NewFindAdminListLogic(r.Context(), svcCtx)
		resp, err := l.FindAdminList(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindAdminListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
