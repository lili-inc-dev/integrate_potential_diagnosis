package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindPersonalityDiagnoseAnswerResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPersonalityDiagnoseAnswerResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPersonalityDiagnoseAnswerResultLogic {
	return &FindPersonalityDiagnoseAnswerResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPersonalityDiagnoseAnswerResultLogic) FindPersonalityDiagnoseAnswerResult() (resp *types.FindPersonalityDiagnoseAnswerResultRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindPersonalityDiagnoseAnswerResult error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return nil, err
	}
	resultList, err := l.svcCtx.PersonalityDiagnoseAnswerModel.FindAnswerResultList(l.ctx, user.Id)
	if err != nil {
		return nil, err
	}

	responseResult := make([]types.PersonalityDiagnoseAnswerResult, len(resultList))
	for i, result := range resultList {
		responseResult[i] = types.PersonalityDiagnoseAnswerResult{
			PersonalityName: result.PersonalityName,
			TotalPoint:      result.TotalPoint,
		}
	}
	return &types.FindPersonalityDiagnoseAnswerResultRes{
		Results: responseResult,
	}, nil
}
