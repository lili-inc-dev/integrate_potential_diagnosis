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

func FrontUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.UserBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserListHandler error"))
			return
		}

		var req types.FrontUserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserListHandler error"))
			return
		}

		l := logic.NewFrontUserListLogic(r.Context(), svcCtx)
		resp, err := l.FrontUserList(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FrontUserListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
