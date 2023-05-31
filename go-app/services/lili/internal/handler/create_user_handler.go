package handler

import (
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/middleware"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetFirebaseTokenFromReq(svcCtx.FirebaseAuth, r)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "CreateUserHandler error"))
			return
		}

		signUpStateClaim, ok := token.Claims[constant.FirebaseCustomClaimKeySignUpState]
		if !ok {
			err := errors.Errorf("firebase id token must not have `%s` claim", constant.FirebaseCustomClaimKeySignUpState)
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "CreateUserHandler error"))
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
			httpx.Error(w, errors.Wrap(err, "CreateUserHandler error"))
			return
		}

		var req types.CreateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateUserHandler error"))
			return
		}

		l := logic.NewCreateUserLogic(r.Context(), svcCtx)
		err = l.CreateUser(&req, token.UID)
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "CreateUserHandler error"))
		} else {
			httpx.Ok(w)
		}
	}
}
