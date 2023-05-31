package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAdminRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindAdminRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAdminRoleListLogic {
	return &FindAdminRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FindAdminRoleList AdminRoleの一覧を取得
func (l *FindAdminRoleListLogic) FindAdminRoleList(ctx context.Context) (resp *types.FindAdminRoleListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAdminRoleList error")
	}()

	adminRoles, err := l.svcCtx.AdminRoleModel.FindAdminRoleWithGrantedCountListAll()
	if err != nil {
		return nil, err
	}

	responseAdminRoles := make([]types.AdminRoleWithGrantedCount, len(adminRoles))
	for i := range adminRoles {
		adminRole := adminRoles[i]

		var description *string
		if adminRole.Description.Valid {
			description = &adminRole.Description.String
		}
		responseAdminRole := types.AdminRoleWithGrantedCount{
			Id:           adminRole.Id,
			Name:         adminRole.Name,
			Description:  description,
			GrantedCount: adminRole.GrantedCount,
		}
		responseAdminRoles[i] = responseAdminRole
	}

	return &types.FindAdminRoleListRes{
		AdminRoles: responseAdminRoles,
	}, nil
}
