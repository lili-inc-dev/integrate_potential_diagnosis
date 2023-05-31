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

func UpdateNoticeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.NoticeEditable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateNoticeHandler error"))
			return
		}

		var req types.UpdateNoticeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateNoticeHandler error"))
			return
		}

		l := logic.NewUpdateNoticeLogic(r.Context(), svcCtx)
		err := l.UpdateNotice(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateNoticeHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
