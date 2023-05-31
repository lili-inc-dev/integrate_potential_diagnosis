package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.AdminEditable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateAdminHandler error"))
			return
		}

		var req types.UpdateAdminReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateAdminHandler error"))
			return
		}

		l := logic.NewUpdateAdminLogic(r.Context(), svcCtx)
		err := l.UpdateAdmin(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateAdminHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
