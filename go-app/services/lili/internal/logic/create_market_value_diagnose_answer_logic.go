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

type CreateMarketValueDiagnoseAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMarketValueDiagnoseAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMarketValueDiagnoseAnswerLogic {
	return &CreateMarketValueDiagnoseAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMarketValueDiagnoseAnswerLogic) CreateMarketValueDiagnoseAnswer(req *types.CreateMarketValueDiagnoseAnswerReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateMarketValueDiagnoseAnswer error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}
	answers := make([]repository.MarketValueDiagnoseAnswers, len(req.Answers))
	for i, reqAnswer := range req.Answers {
		answers[i] = repository.MarketValueDiagnoseAnswers{
			QuestionId: reqAnswer.QuestionId,
			ChoiceId:   reqAnswer.ChoiceId,
			UserId:     user.Id,
		}
	}
	if err := l.svcCtx.MarketValueDiagnoseAnswerModel.BulkInsert(answers); err != nil {
		return err
	}

	csStatus, err := l.svcCtx.UserCsStatusModel.FindOneByUserIdNoCache(l.ctx, user.Id)
	if err != nil {
		return err
	}

	// 診断中ステータスの場合に更新する
	if csStatus.Status == repository.CsStatusDiagnosing {
		PersonalityQuestionCount, err := l.svcCtx.PersonalityDiagnoseQuestionModel.FindCount(l.ctx)
		if err != nil {
			return err
		}

		personalityAnswers, err := l.svcCtx.PersonalityDiagnoseAnswerModel.FindLatestAnswerList(l.ctx, user.Id)
		if err != nil {
			return err
		}
		careerWorkAnswers, err := l.svcCtx.CareerWorkAnswerModel.FetchListByUserId(l.ctx, user.Id)
		if err != nil {
			return err
		}

		var isFinishPersonality bool
		var isFinishCareerWork bool

		if *PersonalityQuestionCount == uint64(len(personalityAnswers)) {
			isFinishPersonality = true
		}
		if len(careerWorkAnswers) > 0 {
			isFinishCareerWork = true
		}

		// 3つとも診断が終わっていたら更新
		if isFinishPersonality && isFinishCareerWork {
			if err := l.svcCtx.UserCsStatusModel.UpdateStatus(l.ctx, csStatus.Id, repository.CsStatusNormal); err != nil {
				return err
			}
		}
	}

	return nil
}
