package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAdminListLogic {
	return &FindAdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAdminListLogic) FindAdminList(req *types.FindAdminListReq) (resp *types.FindAdminListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAdminList error")
	}()

	admins, err := l.svcCtx.AdminModel.FindSpecificRangeListWithRoleNoCache(l.ctx, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		return nil, err
	}
	count, err := l.svcCtx.AdminModel.FindCount(l.ctx)
	if err != nil {
		return nil, err
	}

	firebaseUIDList := []string{}
	for _, ad := range admins {
		firebaseUIDList = append(firebaseUIDList, ad.FirebaseUid)
	}
	fUserList, err := l.svcCtx.FirebaseAuth.FindUserList(l.ctx, firebaseUIDList)
	if err != nil {
		return nil, err
	}

	fUIDToEmail := make(map[string]string)
	for _, fUser := range fUserList {
		fUIDToEmail[fUser.UID] = fUser.Email
	}

	responseAdmin := make([]types.AdminWithRoleInfo, len(admins))
	for i := range admins {
		admin := admins[i]

		var affiliation *string
		if admin.Affiliation.Valid {
			affiliation = &admin.Affiliation.String
		}
		responseAdmin[i] = types.AdminWithRoleInfo{
			Id:           admin.Id,
			RoleId:       admin.RoleId,
			RoleName:     admin.RoleName,
			Email:        fUIDToEmail[admin.FirebaseUid],
			Name:         admin.Name,
			LastAccessAt: util.TimeToStringFormatJp(admin.LastAccessAt),
			Affiliation:  affiliation,
			Status:       admin.IsDisabled,
		}
	}

	pageCount := (*count + req.PageSize - 1) / req.PageSize

	return &types.FindAdminListRes{
		PageCount: pageCount,
		Admins:    responseAdmin,
	}, nil
}
