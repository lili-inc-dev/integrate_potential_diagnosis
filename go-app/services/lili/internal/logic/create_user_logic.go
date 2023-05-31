package logic

import (
	"context"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq, firebaseUid string) (err error) {
	defer func() {
		err = errors.Wrap(err, "CreateUser error")
	}()

	inactiveUser, err := l.svcCtx.InactiveUserModel.FindByFirebaseUidNoCache(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	ulid, err := util.GenerateUlid()
	if err != nil {
		return err
	}
	userId := ulid.String()

	birthDay, err := util.DateStrToTime(req.Birthday)
	if err != nil {
		return err
	}

	if err := l.svcCtx.UserModel.InsertUserNoCache(
		l.ctx,
		userId,
		inactiveUser.LineId,
		inactiveUser.FirebaseUid,
		req.PhoneNumber,
		req.Name,
		inactiveUser.TypeId,
		req.GenderId,
		birthDay,
	); err != nil {
		return err
	}

	if err := l.svcCtx.UserCsStatusModel.UpsertNoCache(
		l.ctx,
		userId,
		repository.CsStatusRegisterOnly,
	); err != nil {
		return err
	}

	ulid, err = util.GenerateUlid()
	if err != nil {
		return err
	}

	if err := l.svcCtx.UserAddressesModel.InsertNoCache(l.ctx, ulid.String(), userId, req.PostalCode, req.Address); err != nil {
		return err
	}

	ulid, err = util.GenerateUlid()
	if err != nil {
		return err
	}

	if err = l.svcCtx.GeneralUserProfileModel.InsertNoCache(
		l.ctx,
		ulid.String(),
		req.NickName,
		req.University,
		req.Faculty,
		req.Department,
		req.GraduationYear,
		req.ServiceTriggerId,
		req.ServiceTriggerDetail,
		req.Introducer,
		req.DesiredAnnualIncomeId,
		userId,
	); err != nil {
		return err
	}

	if len(req.InterestTopicIds) > 0 {
		err = l.svcCtx.UserInterestTopicModel.BulkInsert(userId, req.InterestTopicIds)
		if err != nil {
			return err
		}
	}

	fbUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, firebaseUid)
	if err != nil {
		return err
	}

	// フロントで登録状態を確認するためのカスタムクレイムをセット
	fbUser.CustomClaims[constant.FirebaseCustomClaimKeySignUpState] = constant.FirebaseCustomClaimValueSignUpStateRegisterd
	if err := l.svcCtx.FirebaseAuth.SetCustomClaims(l.ctx, fbUser.UID, fbUser.CustomClaims); err != nil {
		return err
	}

	return nil
}
