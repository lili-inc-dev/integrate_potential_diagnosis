package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCareerWorkAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCareerWorkAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCareerWorkAnswerLogic {
	return &CreateCareerWorkAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCareerWorkAnswerLogic) CreateCareerWorkAnswer(req *types.CreateCareerWorkAnswerReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateCareerWorkAnswer error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}
	answers := make([]repository.CareerWorkAnswers, len(req.Answers))
	for i, reqAnswer := range req.Answers {
		answers[i] = repository.CareerWorkAnswers{
			QuestionKey: reqAnswer.QuestionKey,
			Answer:      reqAnswer.Answer,
			UserId:      user.Id,
			Index:       reqAnswer.Index,
		}
	}
	if err := l.svcCtx.CareerWorkAnswerModel.BulkInsert(answers); err != nil {
		return err
	}

	csStatus, err := l.svcCtx.UserCsStatusModel.FindOneByUserIdNoCache(l.ctx, user.Id)
	if err != nil {
		return err
	}

	if csStatus.Status == repository.CsStatusDiagnosing {
		marketValueQuestionCount, err := l.svcCtx.MarketValueDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return err
		}
		PersonalityQuestionCount, err := l.svcCtx.PersonalityDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return err
		}

		personalityAnswers, err := l.svcCtx.PersonalityDiagnoseAnswerModel.FindLatestAnswerList(l.ctx, user.Id)
		if err != nil {
			return err
		}
		marketValueAnswers, err := l.svcCtx.MarketValueDiagnoseAnswerModel.FindLatestAnswerList(l.ctx, user.Id)
		if err != nil {
			return err
		}
		var isFinishPersonality bool
		var isFinishMarketValue bool

		if *PersonalityQuestionCount == uint64(len(personalityAnswers)) {
			isFinishPersonality = true
		}
		if *marketValueQuestionCount == uint64(len(marketValueAnswers)) {
			isFinishMarketValue = true
		}

		// 3つとも診断が終わっていたら更新
		if isFinishPersonality && isFinishMarketValue {
			if err := l.svcCtx.UserCsStatusModel.UpdateStatus(l.ctx, csStatus.Id, repository.CsStatusNormal); err != nil {
				return err
			}
		}
	}

	return nil
}
