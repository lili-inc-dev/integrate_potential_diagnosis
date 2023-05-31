package logic

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmInactiveUserEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmInactiveUserEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmInactiveUserEmailLogic {
	return &ConfirmInactiveUserEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	emailAuthCodeCheckAttemptCountLimit = 10 // メール認証コード検証の上限試行回数
)

func (l *ConfirmInactiveUserEmailLogic) ConfirmInactiveUserEmail(req *types.ConfirmInactiveUserEmailReq, firebaseUid string) (err error) {
	defer func() {
		err = errors.Wrap(err, "ConfirmInactiveUserEmail error")
	}()

	inactiveUser, err := l.svcCtx.InactiveUserModel.FindByFirebaseUidNoCache(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	emailAuthCode, err := l.svcCtx.EmailAuthenticationCodeModel.FindByInactiveUserIdNoCache(l.ctx, inactiveUser.Id)
	if err != nil {
		return err
	}

	expiration := emailAuthCode.CreatedAt.Add(emailAuthCodeExpirationHour * time.Hour)
	if time.Now().After(expiration) {
		err = errors.New("email auth code has expired")
		err = errorx.New(err, errorx.EmailAuthCodeExpired)
		return err
	}

	if emailAuthCode.AttemptCount >= emailAuthCodeCheckAttemptCountLimit {
		err = errors.New("the attempt count of email auth code check has reached the limit")
		err = errorx.New(err, errorx.ReachedEmailAuthCodeCheckAttemptLimit)
		return err
	}

	if _, err := l.svcCtx.EmailAuthenticationCodeModel.UpdateAttemptCountNoCache(l.ctx, emailAuthCode.Id, emailAuthCode.AttemptCount+1); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(emailAuthCode.CodeHash), []byte(req.AuthCode)); err != nil {
		err = errorx.New(err, errorx.WrongEmailAuthCode)
		return err
	}

	if err := l.svcCtx.FirebaseAuth.UpdateEmailVerified(l.ctx, firebaseUid, true); err != nil {
		return err
	}

	fUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	// フロントで登録状態を確認するためのカスタムクレイムをセット
	fUser.CustomClaims[constant.FirebaseCustomClaimKeySignUpState] = constant.FirebaseCustomClaimValueSignUpStateEmailVerified
	err = l.svcCtx.FirebaseAuth.SetCustomClaims(l.ctx, fUser.UID, fUser.CustomClaims)
	if err != nil {
		return err
	}

	return nil
}
