package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindNoticeListForAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindNoticeListForAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindNoticeListForAdminLogic {
	return &FindNoticeListForAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindNoticeListForAdminLogic) FindNoticeListForAdmin(req *types.FindNoticeListForAdminReq) (resp *types.FindNoticeListForAdminRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindNoticeListForAdmin error")
	}()

	notices, err := l.svcCtx.NoticeModel.FindSpecificRangeNoticeAll(l.ctx, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		return nil, err
	}
	responseNotices := make([]types.Notice, len(notices))
	for i, notice := range notices {
		responseNotices[i] = types.Notice{
			Id:         notice.Id,
			Title:      notice.Title,
			Content:    notice.Content,
			IsReleased: notice.IsReleased,
			CreatedAt:  util.TimeToStringFormatJp(notice.CreatedAt),
			UpdatedAt:  util.TimeToStringFormatJp(notice.UpdatedAt),
		}
	}
	count, err := l.svcCtx.NoticeModel.FindCountNoticeAll(l.ctx)
	if err != nil {
		return nil, err
	}

	pageCount := (*count + req.PageSize - 1) / req.PageSize

	return &types.FindNoticeListForAdminRes{
		Notices:   responseNotices,
		PageCount: pageCount,
	}, nil
}
