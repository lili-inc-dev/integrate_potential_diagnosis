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

func UpdateFrontUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.UserEditable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateFrontUserHandler error"))
			return
		}

		var req types.UpdateFrontUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateFrontUserHandler error"))
			return
		}

		l := logic.NewUpdateFrontUserLogic(r.Context(), svcCtx)
		err := l.UpdateFrontUser(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateFrontUserHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
