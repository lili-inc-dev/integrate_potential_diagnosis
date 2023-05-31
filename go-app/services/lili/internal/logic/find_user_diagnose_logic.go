package logic

import (
	"context"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserDiagnoseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindUserDiagnoseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserDiagnoseLogic {
	return &FindUserDiagnoseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindUserDiagnoseLogic) FindUserDiagnose(req *types.DiagnoseResultReq) (resp *types.DiagnoseResultRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserDiagnose error")
	}()

	if req.DiagnoseType == "personality" {
		personalityAnswerResults, err := l.svcCtx.PersonalityDiagnoseAnswerModel.FindAnswerResultListByAnswerGroupId(l.ctx, req.AnswerGroupId)
		if err != nil {
			return nil, err
		}

		responsePersonalityAnswerResults := make([]types.PersonalityDiagnoseAnswerResult, len(personalityAnswerResults))
		for i, personalityAnswerResult := range personalityAnswerResults {
			responsePersonalityAnswerResults[i] = types.PersonalityDiagnoseAnswerResult{
				PersonalityName: personalityAnswerResult.PersonalityName,
				TotalPoint:      personalityAnswerResult.TotalPoint,
			}
		}

		return &types.DiagnoseResultRes{
			PersonalityDiagnoseAnswerResult: responsePersonalityAnswerResults,
		}, nil
	} else if req.DiagnoseType == "market_value" {
		marketValueResults, err := l.svcCtx.MarketValueDiagnoseAnswerModel.FindAnswerResultListByAnswerGroupId(l.ctx, req.AnswerGroupId)
		if err != nil {
			return nil, err
		}

		responseMarketValueResults := make([]types.MarketValueDiagnoseAnswerResult, len(marketValueResults))
		for i, marketValueResult := range marketValueResults {
			responseMarketValueResults[i] = types.MarketValueDiagnoseAnswerResult{
				MarketValueName: marketValueResult.MarketValueName,
				TotalPoint:      marketValueResult.TotalPoint,
			}
		}

		return &types.DiagnoseResultRes{
			MarketValueDiagnoseAnswerResult: responseMarketValueResults,
		}, nil
	} else if req.DiagnoseType == "career_work" {
		careerWorkAnswers, err := l.svcCtx.CareerWorkAnswerModel.FindListByAnswerGroupId(l.ctx, req.AnswerGroupId)
		if err != nil {
			return nil, err
		}

		responseCareerWorkAnswers := make([]types.CareerWorkAnswer, len(careerWorkAnswers))
		for i, careerWorkAnswer := range careerWorkAnswers {
			responseCareerWorkAnswers[i] = types.CareerWorkAnswer{
				QuestionKey: careerWorkAnswer.QuestionKey,
				Answer:      careerWorkAnswer.Answer,
				Index:       careerWorkAnswer.Index,
			}
		}

		return &types.DiagnoseResultRes{
			CareerWorkAnswer: responseCareerWorkAnswers,
		}, nil
	}

	return nil, errors.New("not match diagnose type")
}
