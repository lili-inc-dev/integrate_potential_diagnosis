package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindMyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindMyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMyPageLogic {
	return &FindMyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindMyPageLogic) FindMyPage() (resp *types.FindMyPageRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindMyPage error")
	}()

	admin, ok := l.ctx.Value(constant.CtxKeyAdmin).(*repository.Admins)
	if !ok {
		err = errors.New("context cast error")
		return nil, err
	}

	fUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, admin.FirebaseUid)
	if err != nil {
		return nil, err
	}

	return &types.FindMyPageRes{
		Email: fUser.Email,
	}, nil
}
