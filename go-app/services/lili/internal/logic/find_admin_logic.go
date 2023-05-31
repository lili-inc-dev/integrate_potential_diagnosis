package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAdminLogic {
	return &FindAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAdminLogic) FindAdmin(req *types.FindAdminReq) (resp *types.FindAdminRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAdmin error")
	}()

	admin, err := l.svcCtx.AdminModel.FindOneWithRoleNoCache(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	var affiliation *string
	if admin.Affiliation.Valid {
		affiliation = &admin.Affiliation.String
	}

	fUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, admin.FirebaseUid)
	if err != nil {
		return nil, err
	}

	return &types.FindAdminRes{
		Admin: types.AdminWithRoleInfo{
			Id:           admin.Id,
			RoleId:       admin.RoleId,
			RoleName:     admin.RoleName,
			Email:        fUser.Email,
			Name:         admin.Name,
			LastAccessAt: util.TimeToStringFormatJp(admin.LastAccessAt),
			Affiliation:  affiliation,
			Status:       admin.IsDisabled,
		},
	}, nil
}
