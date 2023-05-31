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

type UpdateMyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyPageLogic {
	return &UpdateMyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyPageLogic) UpdateMyPage(req *types.UpdateMyPageReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateMyPage error")
	}()

	admin, ok := l.ctx.Value(constant.CtxKeyAdmin).(*repository.Admins)
	if !ok {
		err := errors.New("context cast error")
		return err
	}
	if err := l.svcCtx.FirebaseAuth.ChangePassword(l.ctx, admin.FirebaseUid, req.Password); err != nil {
		return err
	}
	return nil
}
