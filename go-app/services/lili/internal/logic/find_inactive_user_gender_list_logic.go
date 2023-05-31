package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInactiveUserGenderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInactiveUserGenderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInactiveUserGenderListLogic {
	return &FindInactiveUserGenderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInactiveUserGenderListLogic) FindInactiveUserGenderList() (resp *types.FindGenderListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInactiveUserGenderList error")
	}()

	_l := NewFindGenderListLogic(l.ctx, l.svcCtx)
	return _l.FindGenderList()
}
