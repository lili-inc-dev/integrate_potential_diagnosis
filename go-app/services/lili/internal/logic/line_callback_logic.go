package logic

import (
	"context"

	"github.com/pkg/errors"

	"firebase.google.com/go/v4/auth"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/external"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type LineCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLineCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LineCallbackLogic {
	return &LineCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LineCallbackLogic) LineCallback(authCode string) (resp *types.LineCallbackRes, err error) {
	defer func() {
		err = errors.Wrap(err, "LineCallback error")
	}()

	accessToken, err := l.svcCtx.LineAPI.FetchAccessToken(l.ctx, authCode)
	if err != nil {
		return nil, err
	}

	profile, err := l.svcCtx.LineAPI.FetchLineProfileByAccessToken(l.ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if _, err = l.svcCtx.LineAccountModel.Upsert(
		l.ctx,
		profile.LineID,
		profile.DisplayName,
		profile.IconURL,
		profile.StatusMessage,
	); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindByLineIdExceptUnregisteredNoCache(l.ctx, profile.LineID)
	isNotFound := errors.Is(err, repository.ErrNotFound)
	if err != nil && !isNotFound {
		return nil, err
	}

	var fbUser *auth.UserRecord

	if isNotFound {
		// 未登録の場合
		fbUser, err = l.createFirebaseUserForSignUpNotYet(profile)
		if err != nil {
			return nil, err
		}
	} else {
		// 本登録済の場合
		fbUser, err = l.svcCtx.FirebaseAuth.FindUser(l.ctx, user.FirebaseUid)
		if err != nil {
			return nil, err
		}
	}

	// firebase uid と line id を紐付けるためのカスタムクレイムをセット
	// 仮登録が済むまで間紐づける役割
	fbUser.CustomClaims[constant.FirebaseCustomClaimKeyLineID] = profile.LineID
	if err := l.svcCtx.FirebaseAuth.SetCustomClaims(l.ctx, fbUser.UID, fbUser.CustomClaims); err != nil {
		return nil, err
	}

	token, err := l.svcCtx.FirebaseAuth.CreateCustomToken(l.ctx, fbUser.UID)
	if err != nil {
		return nil, err
	}

	return &types.LineCallbackRes{
		FirebaseCustomToken: token,
	}, nil
}

// 未本登録のLINEユーザーに対しfirebase userを作成する
// ただし仮登録途中の場合はfirebase userを使い回す
func (l *LineCallbackLogic) createFirebaseUserForSignUpNotYet(profile *external.LineProfile) (fbUser *auth.UserRecord, err error) {
	defer func() {
		err = errors.Wrap(err, "createFirebaseUserForSignUpNotYet error")
	}()

	inactiveUser, err := l.svcCtx.InactiveUserModel.FindByLineIdNoCache(l.ctx, profile.LineID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, repository.ErrNotFound) {
		// 未仮登録の場合
		fbUser, err := l.svcCtx.FirebaseAuth.CreateUser(l.ctx, profile.DisplayName)
		if err != nil {
			return nil, err
		}

		return fbUser, nil
	}

	oldFbUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, inactiveUser.FirebaseUid)
	if err != nil && !errorx.HasAppErrCode(err, errorx.ResourceNotFound) {
		return nil, err
	}

	if errorx.HasAppErrCode(err, errorx.ResourceNotFound) {
		// 退会済の場合
		fbUser, err := l.svcCtx.FirebaseAuth.CreateUser(l.ctx, profile.DisplayName)
		if err != nil {
			return nil, err
		}

		return fbUser, nil
	}

	signUpState := oldFbUser.CustomClaims[constant.FirebaseCustomClaimKeySignUpState]
	if signUpState == constant.FirebaseCustomClaimValueSignUpStateRegisterd {
		return nil, errors.New("user data inconsistency happened")
	}

	return oldFbUser, nil
}
