package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInterestTopicListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInterestTopicListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInterestTopicListLogic {
	return &FindInterestTopicListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInterestTopicListLogic) FindInterestTopicList() (resp *types.FindInterestTopicListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInterestTopicList error")
	}()

	interTopics, err := l.svcCtx.InterestTopicModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}

	responseInterTopics := make([]types.InterestTopic, len(interTopics))
	for i, interTopic := range interTopics {
		responseInterTopics[i] = types.InterestTopic{
			Id:   interTopic.Id,
			Name: interTopic.Name,
		}
	}

	return &types.FindInterestTopicListRes{
		InterestTopics: responseInterTopics,
	}, nil
}
