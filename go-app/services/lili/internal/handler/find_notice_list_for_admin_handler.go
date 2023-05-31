package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindNoticeListForAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.NoticeBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeListForAdminHandler error"))
			return
		}

		var req types.FindNoticeListForAdminReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeListForAdminHandler error"))
			return
		}

		l := logic.NewFindNoticeListForAdminLogic(r.Context(), svcCtx)
		resp, err := l.FindNoticeListForAdmin(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeListForAdminHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
