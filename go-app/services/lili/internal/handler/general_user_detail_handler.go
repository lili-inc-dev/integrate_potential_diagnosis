package handler

import (
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GeneralUserDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkType := func(userType *repository.UserTypes) bool {
			return userType.Name == repository.UserTypeNameGeneral
		}
		if err := authorizeUser(svcCtx.UserTypeModel, r, checkType); err != nil {
			httpx.Error(w, errors.Wrap(err, "GeneralUserDetailHandler error"))
			return
		}
		l := logic.NewGeneralUserDetailLogic(r.Context(), svcCtx)
		resp, err := l.GeneralUserDetail()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "GeneralUserDetailHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
