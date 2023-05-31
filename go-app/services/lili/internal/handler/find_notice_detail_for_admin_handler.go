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

func FindNoticeDetailForAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.NoticeBrowsable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeDetailForAdminHandler error"))
			return
		}

		var req types.FindNoticeDetailForAdminReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeDetailForAdminHandler error"))
			return
		}

		l := logic.NewFindNoticeDetailForAdminLogic(r.Context(), svcCtx)
		resp, err := l.FindNoticeDetailForAdmin(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindNoticeDetailForAdminHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
