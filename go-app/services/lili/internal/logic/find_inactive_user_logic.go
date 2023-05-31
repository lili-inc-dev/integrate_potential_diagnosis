package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindInactiveUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindInactiveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInactiveUserLogic {
	return &FindInactiveUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindInactiveUserLogic) FindInactiveUser(firebaseUid string) (resp *types.FindInactiveUserRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindInactiveUser error")
	}()

	inactiveUser, err := l.svcCtx.InactiveUserModel.FindWithLineInfoByFirebaseUidNoCache(l.ctx, firebaseUid)
	if err != nil {
		return nil, err
	}

	fbUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, firebaseUid)
	if err != nil {
		return nil, err
	}

	return &types.FindInactiveUserRes{
		LineId:        inactiveUser.LineId,
		LineName:      inactiveUser.LineName,
		IconUrl:       inactiveUser.IconUrl.String,
		StatusMessage: inactiveUser.StatusMessage.String,
		TypeId:        inactiveUser.TypeId,
		Email:         fbUser.Email,
		Name:          inactiveUser.Name,
	}, nil
}
