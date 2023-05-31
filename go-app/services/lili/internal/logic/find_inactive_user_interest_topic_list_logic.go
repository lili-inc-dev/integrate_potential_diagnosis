package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInactiveUserInterestTopicListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInactiveUserInterestTopicListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInactiveUserInterestTopicListLogic {
	return &FindInactiveUserInterestTopicListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInactiveUserInterestTopicListLogic) FindInactiveUserInterestTopicList() (resp *types.FindInterestTopicListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInactiveUserInterestTopicList error")
	}()

	_l := NewFindInterestTopicListLogic(l.ctx, l.svcCtx)
	return _l.FindInterestTopicList()
}
