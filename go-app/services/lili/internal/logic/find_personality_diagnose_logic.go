package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindPersonalityDiagnoseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPersonalityDiagnoseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPersonalityDiagnoseLogic {
	return &FindPersonalityDiagnoseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPersonalityDiagnoseLogic) FindPersonalityDiagnose(req *types.FindPersonalityDiagnoseReq) (resp *types.FindPersonalityDiagnoseRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindPersonalityDiagnose error")
	}()

	resp = &types.FindPersonalityDiagnoseRes{}

	if req.QuestionPage != nil {
		questions, err := l.svcCtx.PersonalityDiagnoseQuestionModel.FindSpecificRangeList(l.ctx, req.QuestionPageSize, (*req.QuestionPage-1)*req.QuestionPageSize)
		if err != nil {
			return nil, err
		}

		responseQuestions := make([]types.PersonalityDiagnoseQuestion, len(questions))
		for i, question := range questions {
			responseQuestions[i] = types.PersonalityDiagnoseQuestion{
				Id:            question.Id,
				PersonalityID: question.PersonalityId,
				Index:         question.Index,
				Content:       question.Content,
			}
		}

		resp.Questions = responseQuestions
	}

	if req.TotalQuestionCount {
		count, err := l.svcCtx.PersonalityDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return nil, err
		}

		resp.TotalQuestionCount = count
	}

	if req.Choice {
		choices, err := l.svcCtx.PersonalityDiagnoseChoiceModel.FindAll(l.ctx)
		if err != nil {
			return nil, err
		}

		responseChoices := make([]types.PersonalityDiagnoseChoice, len(choices))
		for i, choice := range choices {
			responseChoices[i] = types.PersonalityDiagnoseChoice{
				Id:   choice.Id,
				Name: choice.Name,
			}
		}

		resp.Choices = responseChoices
	}

	if req.PerfectPoint {
		perfectPoints, err := l.svcCtx.PersonalityDiagnoseQuestionModel.FindPerfectPointList(l.ctx)
		if err != nil {
			return nil, err
		}
		if len(perfectPoints) <= 0 {
			return nil, errors.New("perfect point is not found")
		}

		perfectPoint := perfectPoints[0].PerfectPoint
		for i := 1; i < len(perfectPoints); i++ {
			if perfectPoints[i].PerfectPoint != perfectPoint {
				return nil, errors.New("perfect points must be same")
			}
		}

		resp.PerfectPoint = &perfectPoint
	}

	return resp, nil
}
