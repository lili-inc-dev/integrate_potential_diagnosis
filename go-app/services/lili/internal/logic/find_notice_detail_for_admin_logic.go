package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindNoticeDetailForAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindNoticeDetailForAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindNoticeDetailForAdminLogic {
	return &FindNoticeDetailForAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindNoticeDetailForAdminLogic) FindNoticeDetailForAdmin(req *types.FindNoticeDetailForAdminReq) (resp *types.FindNoticeDetailForAdminRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindNoticeDetailForAdmin error")
	}()

	notice, err := l.svcCtx.NoticeModel.FindOneByIdNoCache(req.Id)
	if err != nil {
		return nil, err
	}
	return &types.FindNoticeDetailForAdminRes{
		Notice: types.Notice{
			Id:         notice.Id,
			Title:      notice.Title,
			Content:    notice.Content,
			IsReleased: notice.IsReleased,
			CreatedAt:  util.TimeToStringFormatJp(notice.CreatedAt),
			UpdatedAt:  util.TimeToStringFormatJp(notice.UpdatedAt),
		},
	}, nil
}
