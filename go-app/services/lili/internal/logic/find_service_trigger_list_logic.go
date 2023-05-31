package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindServiceTriggerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindServiceTriggerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindServiceTriggerListLogic {
	return &FindServiceTriggerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindServiceTriggerListLogic) FindServiceTriggerList() (resp *types.FindServiceTriggerListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindServiceTriggerList error")
	}()

	serviceTriggers, err := l.svcCtx.ServiceTriggerModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}

	responseServiceTriggers := make([]types.ServiceTrigger, len(serviceTriggers))
	for i, serviceTrigger := range serviceTriggers {
		responseServiceTriggers[i] = types.ServiceTrigger{
			Id:   serviceTrigger.Id,
			Name: serviceTrigger.Name,
		}
	}

	return &types.FindServiceTriggerListRes{
		ServiceTriggers: responseServiceTriggers,
	}, nil
}
