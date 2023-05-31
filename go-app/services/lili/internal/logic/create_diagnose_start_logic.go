package logic

import (
	"context"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDiagnoseStartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDiagnoseStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDiagnoseStartLogic {
	return &CreateDiagnoseStartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDiagnoseStartLogic) CreateDiagnoseStart() (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateDiagnoseStart error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}

	csStatus, err := l.svcCtx.UserCsStatusModel.FindOneByUserIdNoCache(l.ctx, user.Id)
	if err != nil {
		return err
	}

	// `登録のみ`ステータスの場合に更新する
	if csStatus.Status == repository.CsStatusRegisterOnly {
		err := l.svcCtx.UserCsStatusModel.UpdateStatus(l.ctx, csStatus.Id, repository.CsStatusDiagnosing)
		if err != nil {
			return err
		}
	}

	return nil
}
