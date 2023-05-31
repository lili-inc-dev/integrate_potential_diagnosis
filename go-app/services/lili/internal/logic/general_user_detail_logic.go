package logic

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeneralUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneralUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneralUserDetailLogic {
	return &GeneralUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneralUserDetailLogic) GeneralUserDetail() (resp *types.GeneralUserDetailRes, err error) {
	defer func() {
		err = errors.Wrap(err, "GeneralUserDetail error")
	}()

	user, ok := l.ctx.Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		err = errors.New("context cast error")
		return nil, err
	}
	userDetail, err := l.svcCtx.UserModel.FindUserDetailByIdNoCache(l.ctx, user.Id)
	if err != nil {
		return nil, err
	}
	birthDay := util.TimeToStringFormatDateJp(userDetail.Birthday)

	interestTopics, err := l.svcCtx.UserInterestTopicModel.FindListByUserId(l.ctx, user.Id)
	if err != nil {
		return nil, err
	}
	responseInterestTopics := make([]types.InterestTopic, len(interestTopics))

	for i, interestTopic := range interestTopics {
		responseInterestTopics[i] = types.InterestTopic{
			Id:   interestTopic.TopicId,
			Name: interestTopic.Name,
		}
	}

	firebaseUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, userDetail.FirebaseUid)
	if err != nil {
		return nil, err
	}

	var serviceTriggerDetail *string
	if userDetail.ServiceTriggerDetail.Valid {
		serviceTriggerDetail = &userDetail.ServiceTriggerDetail.String
	}
	var introducer *string
	if userDetail.Introducer.Valid {
		introducer = &userDetail.Introducer.String
	}

	return &types.GeneralUserDetailRes{
		Email:                    firebaseUser.Email,
		Name:                     userDetail.Name,
		Nickname:                 userDetail.Nickname,
		GenderName:               userDetail.GenderName,
		University:               userDetail.University,
		Faculty:                  userDetail.Faculty,
		Department:               userDetail.Department,
		GraduationYear:           userDetail.GraduationYear,
		Birthday:                 birthDay,
		PhoneNumber:              userDetail.PhoneNumber,
		PostalCode:               userDetail.PostalCode,
		Address:                  userDetail.Address,
		ServiceTriggerId:         userDetail.ServiceTriggerId,
		ServiceTriggerName:       userDetail.ServiceTriggerName,
		ServiceTriggerDetail:     serviceTriggerDetail,
		DesiredAnnualIncomeId:    userDetail.DesiredAnnualIncomeId,
		DesiredAnnualIncomeValue: userDetail.DesiredAnnualIncomeValue,
		InterestTopics:           responseInterestTopics,
		Introducer:               introducer,
		LineId:                   userDetail.LineId,
		LineName:                 userDetail.LineName,
		IconUrl:                  userDetail.IconUrl.String,
		StatusMessage:            userDetail.StatusMessage.String,
	}, nil
}
