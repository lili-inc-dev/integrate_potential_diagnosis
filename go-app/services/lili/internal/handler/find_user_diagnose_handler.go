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

func FindUserDiagnoseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.UserBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindUserDiagnoseHandler error"))
			return
		}

		var req types.DiagnoseResultReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindUserDiagnoseHandler error"))
			return
		}

		l := logic.NewFindUserDiagnoseLogic(r.Context(), svcCtx)
		resp, err := l.FindUserDiagnose(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindUserDiagnoseHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
