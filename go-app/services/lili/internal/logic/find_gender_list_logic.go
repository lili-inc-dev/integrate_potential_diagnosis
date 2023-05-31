package logic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindGenderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindGenderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindGenderListLogic {
	return &FindGenderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindGenderListLogic) FindGenderList() (resp *types.FindGenderListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindGenderList error")
	}()

	genders, err := l.svcCtx.GenderModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}

	responseGenders := make([]types.Gender, len(genders))
	for i, gender := range genders {
		responseGenders[i] = types.Gender{
			Id:   gender.Id,
			Name: gender.Name,
		}
	}

	return &types.FindGenderListRes{
		Genders: responseGenders,
	}, nil
}
