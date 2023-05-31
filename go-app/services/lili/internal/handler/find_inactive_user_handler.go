package handler

import (
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/middleware"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindInactiveUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetFirebaseTokenFromReq(svcCtx.FirebaseAuth, r)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserHandler error"))
			return
		}

		signUpStateClaim, ok := token.Claims[constant.FirebaseCustomClaimKeySignUpState]
		if !ok {
			err := errors.Errorf("firebase id token must have `%s` claim", constant.FirebaseCustomClaimKeySignUpState)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserHandler error"))
			return
		}

		signUpState, ok := signUpStateClaim.(string)
		if !(ok && signUpState == constant.FirebaseCustomClaimValueSignUpStateEmailVerified) {
			err := errors.Errorf(
				"`%s` claim value must be %s",
				constant.FirebaseCustomClaimKeySignUpState,
				constant.FirebaseCustomClaimValueSignUpStateEmailVerified,
			)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserHandler error"))
			return
		}

		l := logic.NewFindInactiveUserLogic(r.Context(), svcCtx)
		resp, err := l.FindInactiveUser(token.UID)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
