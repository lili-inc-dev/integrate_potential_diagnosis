package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindMarketValueDiagnoseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindMarketValueDiagnoseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindMarketValueDiagnoseLogic {
	return &FindMarketValueDiagnoseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindMarketValueDiagnoseLogic) FindMarketValueDiagnose(req *types.FindMarketValueDiagnoseReq) (resp *types.FindMarketValueDiagnoseRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindMarketValueDiagnose error")
	}()

	resp = &types.FindMarketValueDiagnoseRes{}

	if req.QuestionPage != nil {
		questions, err := l.svcCtx.MarketValueDiagnoseQuestionModel.FindListByDisplayOrder(l.ctx, *req.QuestionPage)
		if err != nil {
			return nil, err
		}

		marketValue, err := l.svcCtx.MarketValueModel.FindOneByDisplayOrder(l.ctx, *req.QuestionPage)
		if err != nil {
			return nil, err
		}

		responseQuestions := make([]types.MarketValueDiagnoseQuestion, len(questions))
		for i, question := range questions {
			responseQuestions[i] = types.MarketValueDiagnoseQuestion{
				Id:            question.Id,
				MarketValueId: question.MarketValueId,
				Index:         question.Index,
				Content:       question.Content,
			}
		}

		resp.Questions = responseQuestions
		resp.MarketValueName = &marketValue.Name
	}

	if req.TotalPageCount {
		count, err := l.svcCtx.MarketValueModel.FindCount(l.ctx)
		if err != nil {
			return nil, err
		}

		resp.TotalPageCount = count
	}

	if req.TotalQuestionCount {
		count, err := l.svcCtx.MarketValueDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return nil, err
		}

		resp.TotalQuestionCount = count
	}

	if req.Choice {
		choices, err := l.svcCtx.MarketValueDiagnoseChoiceModel.FindAll(l.ctx)
		if err != nil {
			return nil, err
		}

		responseChoices := make([]types.MarketValueDiagnoseChoice, len(choices))
		for i, choice := range choices {
			responseChoices[i] = types.MarketValueDiagnoseChoice{
				Id:   choice.Id,
				Name: choice.Name,
			}
		}

		resp.Choices = responseChoices
	}

	if req.PerfectPoint {
		perfectPoint := repository.MarketValueDiagnosePerfectPoint
		resp.PerfectPoint = &perfectPoint
	}

	return resp, nil
}
