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

type CreatePersonalityDiagnoseAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePersonalityDiagnoseAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePersonalityDiagnoseAnswerLogic {
	return &CreatePersonalityDiagnoseAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePersonalityDiagnoseAnswerLogic) CreatePersonalityDiagnoseAnswer(req *types.CreatePersonalityDiagnoseAnswerReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreatePersonalityDiagnoseAnswer error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}

	answers := make([]repository.PersonalityDiagnoseAnswers, len(req.Answers))
	for i, reqAnswer := range req.Answers {
		answers[i] = repository.PersonalityDiagnoseAnswers{
			QuestionId: reqAnswer.QuestionId,
			ChoiceId:   reqAnswer.ChoiceId,
			UserId:     user.Id,
		}
	}
	if err := l.svcCtx.PersonalityDiagnoseAnswerModel.BulkInsert(answers); err != nil {
		return err
	}

	csStatus, err := l.svcCtx.UserCsStatusModel.FindOneByUserIdNoCache(l.ctx, user.Id)
	if err != nil {
		return err
	}

	// 診断中ステータスの場合に更新する
	if csStatus.Status == repository.CsStatusDiagnosing {
		marketValueQuestionCount, err := l.svcCtx.MarketValueDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return err
		}

		marketValueAnswers, err := l.svcCtx.MarketValueDiagnoseAnswerModel.FindLatestAnswerList(l.ctx, user.Id)
		if err != nil {
			return err
		}
		careerWorkAnswers, err := l.svcCtx.CareerWorkAnswerModel.FetchListByUserId(l.ctx, user.Id)
		if err != nil {
			return err
		}

		var isFinishMarketValue bool
		var isFinishCareerWork bool

		if *marketValueQuestionCount == uint64(len(marketValueAnswers)) {
			isFinishMarketValue = true
		}
		if len(careerWorkAnswers) > 0 {
			isFinishCareerWork = true
		}

		// 3つとも診断が終わっていたら更新
		if isFinishCareerWork && isFinishMarketValue {
			if err := l.svcCtx.UserCsStatusModel.UpdateStatus(l.ctx, csStatus.Id, repository.CsStatusNormal); err != nil {
				return err
			}
		}
	}

	return nil
}
