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

func CreateNoticeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.NoticeEditable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateNoticeHandler error"))
			return
		}

		var req types.CreateNoticeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateNoticeHandler error"))
			return
		}

		l := logic.NewCreateNoticeLogic(r.Context(), svcCtx)
		err := l.CreateNotice(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateNoticeHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
