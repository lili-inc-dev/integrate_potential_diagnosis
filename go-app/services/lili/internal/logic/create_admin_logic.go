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

type CreateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(req *types.CreateAdminReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateAdmin error")
	}()

	fUser, err := l.svcCtx.FirebaseAuth.CreateUserWithEmailAndPassword(l.ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	var aff sql.NullString
	if err := aff.Scan(req.Affiliation); err != nil {
		return err
	}

	_, err = l.svcCtx.AdminModel.InsertNoCache(l.ctx, &repository.Admins{
		RoleId:      req.RoleId,
		FirebaseUid: fUser.UID,
		Name:        req.Name,
		Affiliation: aff,
		IsDisabled:  req.IsDisabled,
	})
	if err != nil {
		return err
	}
	return nil
}
