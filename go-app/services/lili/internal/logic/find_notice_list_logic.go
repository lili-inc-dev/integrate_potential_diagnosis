package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindNoticeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindNoticeListLogic {
	return &FindNoticeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindNoticeListLogic) FindNoticeList(req *types.FindNoticeListReq) (resp *types.FindNoticeListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindNoticeList error")
	}()

	notices, err := l.svcCtx.NoticeModel.FindSpecificRangeNotice(l.ctx, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		return nil, err
	}
	responseNotices := make([]types.NoticeForUser, len(notices))
	for i, notice := range notices {
		responseNotices[i] = types.NoticeForUser{
			Id:        notice.Id,
			Title:     notice.Title,
			Content:   notice.Content,
			CreatedAt: util.TimeToStringFormatJp(notice.CreatedAt),
			UpdatedAt: util.TimeToStringFormatJp(notice.UpdatedAt),
		}
	}

	count, err := l.svcCtx.NoticeModel.FindCountReleasedNotice(l.ctx)
	if err != nil {
		return nil, err
	}

	pageCount := (*count + req.PageSize - 1) / req.PageSize

	return &types.FindNoticeListRes{
		Notices:   responseNotices,
		PageCount: pageCount,
	}, nil
}
