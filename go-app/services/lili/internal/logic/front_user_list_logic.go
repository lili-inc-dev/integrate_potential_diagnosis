package logic

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type FrontUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFrontUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FrontUserListLogic {
	return &FrontUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FrontUserListLogic) FrontUserList(req *types.FrontUserListReq) (resp *types.FrontUserListRes, err error) {
	defer func() {
		err = errors.Wrap(err, "FrontUserList error")
	}()

	offset := (req.Page - 1) * req.PageSize

	users, totalCount, err := l.fetchUserList(req.Search, req.CsStatus, req.PageSize, offset)

	if err != nil {
		return nil, err
	}

	responseUsers := make([]types.FrontUser, len(users))
	for i := range users {
		user := users[i]

		createdAt := util.TimeToStringFormatJp(user.CreatedAt)

		var lineId, iconUrl string
		if user.Status == repository.UserStatusUnregistered {
			lineId = user.LineIdUnregistered.String
			iconUrl = user.IconUrlUnregistered.String
		} else {
			lineId = user.LineId.String
			iconUrl = user.IconUrl.String
		}

		var memo *string
		if user.Memo.Valid {
			memo = &user.Memo.String
		}

		responseUsers[i] = types.FrontUser{
			Id:        user.Id,
			IconUrl:   iconUrl,
			Name:      user.Name,
			UserType:  user.TypeName,
			LineId:    lineId,
			Memo:      memo,
			Status:    user.Status,
			CsStatus:  user.CsStatus,
			CreatedAt: createdAt,
		}
	}

	pageCount := (totalCount + req.PageSize - 1) / req.PageSize

	return &types.FrontUserListRes{
		PageCount: pageCount,
		Users:     responseUsers,
	}, nil
}

func (l *FrontUserListLogic) fetchUserList(searchKeyword *string, csStatus *string, limit uint64, offset uint64) (users []repository.UserWithLineInfo, tatolCount uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "fetchUserList error")
	}()

	if searchKeyword != nil {
		users, tatolCount, err = l.fetchUserListBySearchKeywordAndCsStatus(*searchKeyword, csStatus, limit, offset)
	} else {
		users, tatolCount, err = l.fetchUserListNoSearchKeyword(csStatus, limit, offset)
	}

	return
}

func (l *FrontUserListLogic) fetchUserListBySearchKeywordAndCsStatus(searchKeyword string, csStatus *string, limit uint64, offset uint64) (resp []repository.UserWithLineInfo, tatolCount uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "fetchUserListBySearchKeywordAndCsStatus error")
	}()

	var users []repository.UserWithLineInfo

	if csStatus != nil {
		users, err = l.svcCtx.UserModel.FindUserWithLineInfoListByCsStatus(l.ctx, *csStatus)
	} else {
		users, err = l.svcCtx.UserModel.FindUserWithLineInfoListAll(l.ctx)
	}

	if err != nil {
		return nil, 0, err
	}

	userIdToEmail := make(map[string]string)
	firebaseUidToUserId := make(map[string]string)
	firebaseUserIds := make([]string, 0)

	for _, user := range users {
		if user.Status == repository.UserStatusUnregistered {
			userIdToEmail[user.Id] = user.EmailUnregistered.String
		} else {
			firebaseUidToUserId[user.FirebaseUid] = user.Id
			firebaseUserIds = append(firebaseUserIds, user.FirebaseUid)
		}
	}

	firebaseUserRecords, err := l.svcCtx.FirebaseAuth.FindUserList(l.ctx, firebaseUserIds)
	if err != nil {
		return nil, 0, err
	}

	for _, fbUser := range firebaseUserRecords {
		userId := firebaseUidToUserId[fbUser.UID]
		userIdToEmail[userId] = fbUser.Email
	}

	matchUsers := make([]repository.UserWithLineInfo, 0)

	count := uint64(0)
	for _, user := range users {
		email := userIdToEmail[user.Id]
		if strings.Contains(email, searchKeyword) || strings.Contains(user.Name, searchKeyword) {
			count++
			if offset <= count && count < offset+limit {
				matchUsers = append(matchUsers, user)
			}
		}
	}
	return matchUsers, count, nil
}

func (l *FrontUserListLogic) fetchUserListNoSearchKeyword(csStatus *string, limit uint64, offset uint64) (users []repository.UserWithLineInfo, totalCount uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "fetchUserListNoSearchKeyword error")
	}()

	if csStatus != nil {
		users, err = l.svcCtx.UserModel.FindSpecificRangeUserWithLineInfoListByCsStatus(l.ctx, *csStatus, limit, offset)
		if err != nil {
			return nil, 0, err
		}

		totalCount, err = l.svcCtx.UserModel.FindCountByCsStatus(l.ctx, *csStatus)
	} else {
		users, err = l.svcCtx.UserModel.FindSpecificRangeUserWithLineInfoList(l.ctx, limit, offset)
		if err != nil {
			return nil, 0, err
		}

		totalCount, err = l.svcCtx.UserModel.FindCount(l.ctx)
	}

	return
}
