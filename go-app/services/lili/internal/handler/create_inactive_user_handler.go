package handler

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/middleware"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateInactiveUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetFirebaseTokenFromReq(svcCtx.FirebaseAuth, r)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
			return
		}

		lineIdClaim, ok := token.Claims[constant.FirebaseCustomClaimKeyLineID]
		if !ok {
			err := errors.Errorf("firebase id token must have `%s` claim", constant.FirebaseCustomClaimKeyLineID)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
			return
		}

		lineID, ok := lineIdClaim.(string)
		if !ok {
			err := errors.Errorf("type mismatch of `%s` claim", constant.FirebaseCustomClaimKeyLineID)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
			return
		}

		_, ok = token.Claims[constant.FirebaseCustomClaimKeySignUpState]
		if ok {
			err := errors.Errorf("firebase id token must not have `%s` claim when creating the inactive user", constant.FirebaseCustomClaimKeySignUpState)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
			return
		}

		var req types.CreateInactiveUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
			return
		}

		l := logic.NewCreateInactiveUserLogic(r.Context(), svcCtx)

		resp, err := l.CreateInactiveUser(&req, token.UID, lineID)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateInactiveUserHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
