package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnregisterUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnregisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnregisterUserLogic {
	return &UnregisterUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnregisterUserLogic) UnregisterUser(req *types.UnregisterUserReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UnregisterUser error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}

	fbUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, user.FirebaseUid)
	if err != nil {
		return err
	}

	ulid, err := util.GenerateUlid()
	if err != nil {
		return err
	}

	// 退会済ユーザーのLINE IDとEmailは、Unique制約により残しておくことができない
	// そのため別テーブルに記録しておく
	if err := l.svcCtx.UnregisteredUserProfileModel.InsertNoCache(l.ctx,
		ulid.String(),
		user.Id,
		user.LineId.String,
		fbUser.Email,
		req.UnregistrationReasonId,
	); err != nil {
		return err
	}

	if err = l.svcCtx.UserModel.UnregisterNoCache(l.ctx, user.Id); err != nil {
		return err
	}

	if err := l.svcCtx.FirebaseAuth.DeleteUser(l.ctx, fbUser.UID); err != nil {
		return err
	}

	return nil
}
