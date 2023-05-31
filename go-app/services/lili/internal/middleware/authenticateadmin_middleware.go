package middleware

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/external"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthenticateAdminMiddleware struct {
	AdminModel   repository.AdminsModel
	FirebaseAuth external.FirebaseAuth
}

func NewAuthenticateAdminMiddleware(
	adminModel repository.AdminsModel,
	firebaseAuth external.FirebaseAuth,
) *AuthenticateAdminMiddleware {
	return &AuthenticateAdminMiddleware{
		AdminModel:   adminModel,
		FirebaseAuth: firebaseAuth,
	}
}

const (
	// 管理画面ログインセッション有効期間（単位：時）
	AdminLoginSessionExpiration = 24 * 3
)

func (m *AuthenticateAdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetFirebaseTokenFromReq(m.FirebaseAuth, r)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "AuthenticateAdminMiddleware error"))
			return
		}

		if err := checkLoginSessionExpiration(r.Context(), m.FirebaseAuth, token); err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "AuthenticateAdminMiddleware error"))
			return
		}

		admin, err := m.AdminModel.FindOneByFirebaseUidNoCache(r.Context(), token.UID)
		if err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "AuthenticateAdminMiddleware error"))
			return
		}

		if admin.IsDisabled {
			err = errors.New("account has been disabled")
			err = errorx.New(err, errorx.Forbidden)
			httpx.Error(w, errors.Wrap(err, "AuthenticateAdminMiddleware error"))
			return
		}

		go func() {
			err := m.AdminModel.UpdateLastAccessAt(context.Background(), admin.Id)
			if err != nil {
				logx.Error(err)
			}
		}()

		ctx := context.WithValue(r.Context(), constant.CtxKeyAdmin, admin)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetFirebaseTokenFromReq(firebaseAuth external.FirebaseAuth, r *http.Request) (token *auth.Token, err error) {
	defer func() {
		err = errors.Wrap(err, "GetFirebaseTokenFromReq error")
	}()

	tokenStr, err := util.GetBearerToken(r)
	if err != nil {
		return nil, err
	}

	token, err = firebaseAuth.VerifyIDToken(r.Context(), tokenStr)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func checkLoginSessionExpiration(ctx context.Context, firebaseAuth external.FirebaseAuth, token *auth.Token) (err error) {
	defer func() {
		err = errors.Wrap(err, "checkLoginSessionExpiration error")
	}()

	// token.AuthTimeはid tokenリフレッシュ時に引き継がれる（リフレッシュしても有効期限は伸びない）
	authTime := time.Unix(token.AuthTime, 0)
	expiration := authTime.Add(time.Hour * AdminLoginSessionExpiration)
	if time.Now().After(expiration) {
		if err := firebaseAuth.RevokeRefreshTokens(ctx, token.UID); err != nil {
			return err
		}
		return errors.New("login session expired")
	}

	return nil
}
