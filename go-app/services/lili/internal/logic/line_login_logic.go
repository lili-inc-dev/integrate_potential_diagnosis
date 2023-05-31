package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LineLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLineLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LineLoginLogic {
	return &LineLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LineLoginLogic) LineLogin() (err error) {
	defer func() {
		err = errors.Wrap(err, "LineLogin error")
	}()

	// todo: add your logic here and delete this line

	return nil
}
