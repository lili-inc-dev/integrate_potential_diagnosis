package handler

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/handler/valid"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateMyPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMyPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateMyPageHandler error"))
			return
		}

		if err := validateUpdateMyPageReq(&req); err != nil {
			err = errorx.New(err, errorx.InvalidParameter)
			httpx.Error(w, errors.Wrap(err, "UpdateMyPageHandler error"))
			return
		}

		l := logic.NewUpdateMyPageLogic(r.Context(), svcCtx)
		err := l.UpdateMyPage(&req)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "UpdateMyPageHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}

func validateUpdateMyPageReq(req *types.UpdateMyPageReq) error {
	if err := valid.ValidatePassword(req.Password); err != nil {
		return err
	}

	return nil
}
