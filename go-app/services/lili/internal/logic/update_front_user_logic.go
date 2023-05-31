package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFrontUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFrontUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFrontUserLogic {
	return &UpdateFrontUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFrontUserLogic) UpdateFrontUser(req *types.UpdateFrontUserReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateFrontUser error")
	}()

	user, err := l.svcCtx.UserModel.FindByIdNoCache(req.Id)
	if err != nil {
		return err
	}
	if user.Status == repository.UserStatusUnregistered {
		err = errors.New("user status is unregistered")
		return
	}

	if err := l.svcCtx.UserModel.UpdateAdminEditParamNoCache(l.ctx, req.Id, req.Memo, req.Status); err != nil {
		return err
	}
	return nil
}
