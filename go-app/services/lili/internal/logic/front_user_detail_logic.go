package logic

import (
	"context"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FrontUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFrontUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FrontUserDetailLogic {
	return &FrontUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FrontUserDetailLogic) FrontUserDetail(req *types.FrontUserDetailReq) (resp *types.FrontUserDetailRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FrontUserDetail error")
	}()

	userDetail, err := l.svcCtx.UserModel.FindUserDetailForAdminByIdNoCache(req.Id)
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

	diagnoseHistories, err := l.svcCtx.DiagnoseModel.FindSpecificRangeDiagnoseHistoryByUserId(req.Id, req.DiagnoseHistoryPageSize, (req.DiagnoseHistoryPage-1)*req.DiagnoseHistoryPageSize)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.DiagnoseModel.FindCountByUserId(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	pageCount := (*count + req.DiagnoseHistoryPageSize - 1) / req.DiagnoseHistoryPageSize

	responseDiagnoseHistories := make([]types.DiagnoseHistory, len(diagnoseHistories))
	for i, diagnoseHistory := range diagnoseHistories {
		responseDiagnoseHistories[i] = types.DiagnoseHistory{
			DiagnoseType:  diagnoseHistory.DiagnoseType,
			DiagnoseName:  diagnoseHistory.DiagnoseName,
			AnswerGroupId: diagnoseHistory.AnswerGroupId,
			CreatedAt:     util.TimeToStringFormatJp(diagnoseHistory.CreatedAt),
		}
	}

	interestTopics, err := l.svcCtx.UserInterestTopicModel.FindListByUserId(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	interestTopicNames := make([]string, len(interestTopics))
	for i, interestTopic := range interestTopics {
		interestTopicNames[i] = interestTopic.Name
	}

	var memo *string
	if userDetail.Memo.Valid {
		memo = &userDetail.Memo.String
	}

	var lineId, iconUrl, email string

	if userDetail.Status == repository.UserStatusUnregistered {
		lineId = userDetail.LineIdUnregistered.String
		iconUrl = userDetail.IconUrlUnregistered.String
		email = userDetail.EmailUnregistered.String
	} else {
		lineId = userDetail.LineId.String
		iconUrl = userDetail.IconUrl.String

		firebaseUser, err := l.svcCtx.FirebaseAuth.FindUser(l.ctx, userDetail.FirebaseUid)
		if err != nil {
			return nil, err
		}

		email = firebaseUser.Email
	}

	return &types.FrontUserDetailRes{
		Id:                       userDetail.Id,
		LineId:                   lineId,
		UserType:                 userDetail.TypeName,
		Gender:                   userDetail.Gender,
		Name:                     userDetail.Name,
		Nickname:                 userDetail.Nickname,
		IconUrl:                  iconUrl,
		University:               userDetail.University,
		Faculty:                  userDetail.Faculty,
		Department:               userDetail.Department,
		GraduationYear:           userDetail.GraduationYear,
		Birthday:                 util.TimeToStringFormatDateJp(userDetail.Birthday),
		Email:                    email,
		PhoneNumber:              userDetail.PhoneNumber,
		PostalCode:               userDetail.PostalCode,
		Address:                  userDetail.Address,
		CreatedAt:                util.TimeToStringFormatJp(userDetail.CreatedAt),
		ServiceTrigger:           userDetail.ServiceTrigger,
		ServiceTriggerDetail:     serviceTriggerDetail,
		Introducer:               introducer,
		Memo:                     memo,
		Status:                   userDetail.Status,
		CsStatus:                 userDetail.CsStatus,
		DesiredAnnualIncome:      userDetail.DesiredAnnualIncome,
		DiagnoseHistories:        responseDiagnoseHistories,
		DiagnoseHistoryPageCount: pageCount,
		InterestTopics:           interestTopicNames,
	}, nil
}
