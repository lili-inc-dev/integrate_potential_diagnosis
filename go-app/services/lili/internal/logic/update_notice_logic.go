package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoticeLogic {
	return &UpdateNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNoticeLogic) UpdateNotice(req *types.UpdateNoticeReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateNotice error")
	}()

	_, err = l.svcCtx.NoticeModel.FindOneByIdNoCache(req.Id)
	if err != nil {
		return err
	}

	err = l.svcCtx.NoticeModel.UpdateNoCache(l.ctx, req.Id, req.Title, req.Content, req.IsReleased)
	if err != nil {
		return err
	}
	return nil
}
