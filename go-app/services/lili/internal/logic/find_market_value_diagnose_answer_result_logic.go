package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindMarketValueDiagnoseAnswerResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindMarketValueDiagnoseAnswerResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMarketValueDiagnoseAnswerResultLogic {
	return &FindMarketValueDiagnoseAnswerResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindMarketValueDiagnoseAnswerResultLogic) FindMarketValueDiagnoseAnswerResult() (resp *types.FindMarketValueDiagnoseAnswerResultRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindMarketValueDiagnoseAnswerResult error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err = errors.New("context cast error")
		return nil, err
	}
	resultList, err := l.svcCtx.MarketValueDiagnoseAnswerModel.FindAnswerResultList(l.ctx, user.Id)
	if err != nil {
		return nil, err
	}
	responseResult := make([]types.MarketValueDiagnoseAnswerResult, len(resultList))
	for i, result := range resultList {
		responseResult[i] = types.MarketValueDiagnoseAnswerResult{
			MarketValueName: result.MarketValueName,
			TotalPoint:      result.TotalPoint,
		}
	}
	return &types.FindMarketValueDiagnoseAnswerResultRes{
		Results: responseResult,
	}, nil
}
