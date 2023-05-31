package logic

import (
	"context"
	"time"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type ResendAuthCodeRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResendAuthCodeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResendAuthCodeRequestLogic {
	return &ResendAuthCodeRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResendAuthCodeRequestLogic) ResendAuthCodeRequest(firebaseUid string) error {
	authCode, err := generateEmailAuthCode()
	if err != nil {
		return err
	}

	authCodeHash, err := util.GenerateHash(authCode)
	if err != nil {
		return err
	}

	inactiveUser, err := l.svcCtx.InactiveUserModel.FindByFirebaseUidNoCache(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	emailAuthCode, err := l.svcCtx.EmailAuthenticationCodeModel.FindByInactiveUserIdNoCache(l.ctx, inactiveUser.Id)
	if err != nil {
		return err
	}

	expiration := emailAuthCode.CreatedAt.Add(emailAuthCodeExpirationHour * time.Hour)
	isExpired := time.Now().After(expiration)
	if isExpired {
		// 新しい認証コードレコードを発行する

		ulid, err := util.GenerateUlid()
		if err != nil {
			return err
		}

		if _, err = l.svcCtx.EmailAuthenticationCodeModel.InsertNoCache(l.ctx, &repository.EmailAuthenticationCodes{
			Id:             ulid.String(),
			InactiveUserId: inactiveUser.Id,
			CodeHash:       authCodeHash,
			AttemptCount:   0,
		}); err != nil {
			return err
		}
	} else {
		// 既存の認証コードレコードの`code_hash`のみ更新する
		if _, err = l.svcCtx.EmailAuthenticationCodeModel.UpdateCodeHashNoCache(l.ctx, emailAuthCode.Id, authCodeHash); err != nil {
			return err
		}
	}

	fbUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	if err := sendConfirmationMail(l.ctx, l.svcCtx.Email, fbUser.Email, authCode); err != nil {
		return err
	}

	return nil
}
