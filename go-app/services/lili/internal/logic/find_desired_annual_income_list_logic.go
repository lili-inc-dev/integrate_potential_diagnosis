package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindDesiredAnnualIncomeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindDesiredAnnualIncomeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDesiredAnnualIncomeListLogic {
	return &FindDesiredAnnualIncomeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindDesiredAnnualIncomeListLogic) FindDesiredAnnualIncomeList() (resp *types.FindDesiredAnnualIncomeListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindDesiredAnnualIncomeList error")
	}()

	desiredAnnualIncomes, err := l.svcCtx.DesiredAnnualIncomeModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}

	responseDesiredAnnualIncomes := make([]types.DesiredAnnualIncome, len(desiredAnnualIncomes))
	for i, desiredAnnualIncome := range desiredAnnualIncomes {
		responseDesiredAnnualIncomes[i] = types.DesiredAnnualIncome{
			Id:    desiredAnnualIncome.Id,
			Value: desiredAnnualIncome.Value,
		}
	}

	return &types.FindDesiredAnnualIncomeListRes{
		DesiredAnnualIncomes: responseDesiredAnnualIncomes,
	}, nil
}
