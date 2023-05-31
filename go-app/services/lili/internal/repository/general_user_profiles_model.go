package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GeneralUserProfilesModel = (*customGeneralUserProfilesModel)(nil)

type (
	// GeneralUserProfilesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGeneralUserProfilesModel.
	GeneralUserProfilesModel interface {
		generalUserProfilesModel
		FindByUserIdNoCache(userId string) (*GeneralUserProfiles, error)
		FindByUserId(userId string) (*GeneralUserProfileWithUserInfo, error)
		InsertNoCache(
			ctx context.Context,
			id string,
			nickname string,
			university string,
			faculty string,
			department string,
			graduationYear uint64,
			serviceTriggerId uint64,
			serviceTriggerDetail *string,
			introducer *string,
			desiredAnnualIncomeId uint64,
			userId string,
		) error
		UpdateNoCache(
			ctx context.Context,
			nickname string,
			university string,
			faculty string,
			department string,
			graduationYear uint64,
			serviceTriggerId uint64,
			serviceTriggerDetail *string,
			introducer *string,
			desiredAnnualIncomeId uint64,
			userId string,
		) error
	}

	customGeneralUserProfilesModel struct {
		*defaultGeneralUserProfilesModel
	}

	GeneralUserProfileWithUserInfo struct {
		FirebaseUid          string         `db:"firebase_uid"`
		Name                 string         `db:"name"`
		NickName             string         `db:"nickname"`
		GenderId             uint64         `db:"gender_id"`
		Birthday             time.Time      `db:"birthday"`
		PhoneNumber          string         `db:"phone_number"`
		University           string         `db:"university"`
		Faculty              string         `db:"faculty"`
		Department           string         `db:"department"`
		GraduationYear       uint64         `db:"graduation_year"`
		PostalCode           string         `db:"postal_code"`
		Address              string         `db:"address"`
		ServiceTriggerId     uint64         `db:"service_trigger_id"`
		ServiceTriggerDetail sql.NullString `db:"service_trigger_detail"`
		Introducer           sql.NullString `db:"introducer"`
	}
)

// NewGeneralUserProfilesModel returns a model for the database table.
func NewGeneralUserProfilesModel(conn sqlx.SqlConn, c cache.CacheConf) GeneralUserProfilesModel {
	return &customGeneralUserProfilesModel{
		defaultGeneralUserProfilesModel: newGeneralUserProfilesModel(conn, c),
	}
}

func (m *customGeneralUserProfilesModel) FindByUserIdNoCache(userId string) (profiles *GeneralUserProfiles, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByUserIdNoCache error")
	}()

	var resp GeneralUserProfiles
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ?", generalUserProfilesRows, m.table)
	err = m.QueryRowNoCache(&resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customGeneralUserProfilesModel) FindByUserId(userId string) (profile *GeneralUserProfileWithUserInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByUserId error")
	}()

	var resp GeneralUserProfileWithUserInfo

	query := `
SELECT
	u.firebase_uid,
	u.name,
	gup.nickname,
	u.gender_id,
	u.birthday,
	u.phone_number,
	gup.university,
	gup.faculty,
	gup.department,
	gup.graduation_year,
	ua.postal_code,
	ua.address,
	gup.service_trigger_id,
	gup.service_trigger_detail,
	gup.introducer
FROM general_user_profiles AS gup
	INNER JOIN users AS u
	ON u.id = gup.user_id
	INNER JOIN user_addresses AS ua
	ON ua.user_id = gup.user_id 
WHERE gup.user_id = ?
`
	err = m.QueryRowNoCache(&resp, query, userId)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customGeneralUserProfilesModel) InsertNoCache(
	ctx context.Context,
	id string,
	nickname string,
	university string,
	faculty string,
	department string,
	graduationYear uint64,
	serviceTriggerId uint64,
	serviceTriggerDetail *string,
	introducer *string,
	desiredAnnualIncomeId uint64,
	userId string,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()
	query := `
INSERT INTO general_user_profiles(
	id,
	nickname,
	university,
	faculty,
	department,
	graduation_year,
	service_trigger_id,
	service_trigger_detail,
	introducer,
	desired_annual_income_id,
	user_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		id,
		nickname,
		university,
		faculty,
		department,
		graduationYear,
		serviceTriggerId,
		serviceTriggerDetail,
		introducer,
		desiredAnnualIncomeId,
		userId,
	)
	return err
}

func (m *customGeneralUserProfilesModel) UpdateNoCache(
	ctx context.Context,
	nickname string,
	university string,
	faculty string,
	department string,
	graduationYear uint64,
	serviceTriggerId uint64,
	serviceTriggerDetail *string,
	introducer *string,
	desiredAnnualIncomeId uint64,
	userId string,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateNoCache error")
	}()

	query := `
UPDATE
general_user_profiles
SET 
	nickname = ?,
	university = ?,
	faculty = ?,
	department = ?,
	graduation_year = ?,
	service_trigger_id = ?,
	service_trigger_detail = ?,
	introducer = ?,
	desired_annual_income_id = ?
WHERE user_id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		nickname,
		university,
		faculty,
		department,
		graduationYear,
		serviceTriggerId,
		serviceTriggerDetail,
		introducer,
		desiredAnnualIncomeId,
		userId,
	)
	return err
}
