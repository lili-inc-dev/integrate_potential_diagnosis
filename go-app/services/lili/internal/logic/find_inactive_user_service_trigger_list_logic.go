package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInactiveUserServiceTriggerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInactiveUserServiceTriggerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInactiveUserServiceTriggerListLogic {
	return &FindInactiveUserServiceTriggerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInactiveUserServiceTriggerListLogic) FindInactiveUserServiceTriggerList() (resp *types.FindServiceTriggerListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInactiveUserServiceTriggerList error")
	}()

	_l := NewFindServiceTriggerListLogic(l.ctx, l.svcCtx)
	return _l.FindServiceTriggerList()
}
