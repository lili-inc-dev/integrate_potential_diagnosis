package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoticeLogic {
	return &CreateNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNoticeLogic) CreateNotice(req *types.CreateNoticeReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateNotice error")
	}()

	if _, err := l.svcCtx.NoticeModel.InsertNoCache(l.ctx, req.Title, req.Content, req.IsReleased); err != nil {
		return err
	}
	return nil
}
