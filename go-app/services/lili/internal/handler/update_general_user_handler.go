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

func UpdateGeneralUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkType := func(userType *repository.UserTypes) bool {
			return userType.Name == repository.UserTypeNameGeneral
		}
		if err := authorizeUser(svcCtx.UserTypeModel, r, checkType); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateGeneralUserHandler error"))
			return
		}

		var req types.UpdateGeneralUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateGeneralUserHandler error"))
			return
		}

		l := logic.NewUpdateGeneralUserLogic(r.Context(), svcCtx)
		err := l.UpdateGeneralUser(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateGeneralUserHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
