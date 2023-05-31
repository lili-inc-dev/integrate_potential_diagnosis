package middleware

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/external"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthenticateUserMiddleware struct {
	UserModel    repository.UsersModel
	FirebaseAuth external.FirebaseAuth
}

func NewAuthenticateUserMiddleware(
	userModel repository.UsersModel,
	firebaseAuth external.FirebaseAuth,
) *AuthenticateUserMiddleware {
	return &AuthenticateUserMiddleware{
		UserModel:    userModel,
		FirebaseAuth: firebaseAuth,
	}
}

func (m *AuthenticateUserMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetFirebaseTokenFromReq(m.FirebaseAuth, r)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "AuthenticateUserMiddleware error"))
			return
		}

		user, err := m.UserModel.FindOneByFirebaseUidNoCache(r.Context(), token.UID)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "AuthenticateUserMiddleware error"))
			return
		}

		if user.Status != repository.UserStatusRegistered {
			err = errors.New("account has been disabled")
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "AuthenticateUserMiddleware error"))
			return
		}

		go func() {
			err := m.UserModel.UpdateLastAccessAt(context.Background(), user.Id)
			if err != nil {
				logx.Error(err)
			}
		}()

		ctx := context.WithValue(r.Context(), constant.CtxKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
