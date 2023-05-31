package handler

import (
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/handler/valid"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkRole := func(role *repository.AdminRoles) bool {
			return role.AdminEditable
		}
		if err := authorizeAdmin(svcCtx.AdminRoleModel, r, checkRole); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateAdminHandler error"))
			return
		}

		var req types.CreateAdminReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateAdminHandler error"))
			return
		}

		if err := validateCreateAdminReq(&req); err != nil {
			err = errorx.New(err, errorx.InvalidParameter)
			httpx.Error(w, errors.Wrap(err, "CreateAdminHandler error"))
			return
		}

		l := logic.NewCreateAdminLogic(r.Context(), svcCtx)
		err := l.CreateAdmin(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateAdminHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}

func validateCreateAdminReq(req *types.CreateAdminReq) error {
	if err := valid.ValidatePassword(req.Password); err != nil {
		return err
	}

	return nil
}
