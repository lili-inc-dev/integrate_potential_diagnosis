package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInactiveUserDesiredAnnualIncomeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInactiveUserDesiredAnnualIncomeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInactiveUserDesiredAnnualIncomeListLogic {
	return &FindInactiveUserDesiredAnnualIncomeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInactiveUserDesiredAnnualIncomeListLogic) FindInactiveUserDesiredAnnualIncomeList() (resp *types.FindDesiredAnnualIncomeListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInactiveUserDesiredAnnualIncomeList error")
	}()

	_l := NewFindDesiredAnnualIncomeListLogic(l.ctx, l.svcCtx)
	return _l.FindDesiredAnnualIncomeList()
}
