package logic

import (
	"context"
	"database/sql"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(req *types.UpdateAdminReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateAdmin error")
	}()

	var affiliation sql.NullString
	if req.Affiliation != nil {
		affiliation = sql.NullString{
			String: *req.Affiliation,
			Valid:  true,
		}
	} else {
		affiliation = sql.NullString{
			Valid: false,
		}
	}

	if _, err = l.svcCtx.AdminModel.FindOneNoCache(l.ctx, req.Id); err != nil {
		return err
	}

	err = l.svcCtx.AdminModel.Update(l.ctx, &repository.Admins{
		Id:          req.Id,
		RoleId:      req.RoleId,
		Name:        req.Name,
		Affiliation: affiliation,
		IsDisabled:  req.IsDisabled,
	})
	if err != nil {
		return err
	}

	admin, err := l.svcCtx.AdminModel.FindOneNoCache(l.ctx, req.Id)
	if err != nil {
		return err
	}

	if err := l.svcCtx.FirebaseAuth.ChangeEmail(
		l.ctx,
		admin.FirebaseUid,
		req.Email,
	); err != nil {
		return err
	}

	return nil
}
