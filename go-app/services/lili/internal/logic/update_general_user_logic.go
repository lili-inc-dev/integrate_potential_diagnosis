package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGeneralUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGeneralUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGeneralUserLogic {
	return &UpdateGeneralUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGeneralUserLogic) UpdateGeneralUser(req *types.UpdateGeneralUserReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateGeneralUser error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err := errors.New("context cast error")
		return err
	}

	birthDay, err := util.DateStrToTime(req.Birthday)
	if err != nil {
		return err
	}

	err = l.svcCtx.UserModel.UpdateUserEditParamNoCache(l.ctx, user.Id, req.GenderId, birthDay, req.PhoneNumber, req.Name)
	if err != nil {
		return err
	}

	err = l.svcCtx.UserAddressesModel.UpdateNoCache(l.ctx, user.Id, req.PostalCode, req.Address)
	if err != nil {
		return err
	}

	err = l.svcCtx.GeneralUserProfileModel.UpdateNoCache(
		l.ctx,
		req.NickName,
		req.University,
		req.Faculty,
		req.Department,
		req.GraduationYear,
		req.ServiceTriggerId,
		req.ServiceTriggerDetail,
		req.Introducer,
		req.DesiredAnnualIncomeId,
		user.Id,
	)
	if err != nil {
		return err
	}

	err = l.svcCtx.UserInterestTopicModel.DeleteByUserId(user.Id)
	if err != nil {
		return err
	}

	if len(req.InterestTopicIds) > 0 {
		err = l.svcCtx.UserInterestTopicModel.BulkInsert(user.Id, req.InterestTopicIds)
		if err != nil {
			return err
		}
	}
	return nil
}
