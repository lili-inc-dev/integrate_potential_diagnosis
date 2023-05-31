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

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		FindUserDetailForAdminByIdNoCache(id string) (*UserDetailForAdmin, error)
		FindUserWithLineInfoListAll(ctx context.Context) ([]UserWithLineInfo, error)
		FindUserWithLineInfoListByCsStatus(ctx context.Context, csStatus string) ([]UserWithLineInfo, error)
		FindSpecificRangeUserWithLineInfoList(ctx context.Context, limit uint64, offset uint64) ([]UserWithLineInfo, error)
		FindSpecificRangeUserWithLineInfoListByCsStatus(ctx context.Context, csStatus string, limit uint64, offset uint64) ([]UserWithLineInfo, error)
		FindByIdNoCache(userId string) (*Users, error)
		FindOneByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (*Users, error)
		FindByLineIdExceptUnregisteredNoCache(ctx context.Context, lineId string) (user *Users, err error)
		FindCount(ctx context.Context) (uint64, error)
		FindCountByCsStatus(ctx context.Context, csStatus string) (uint64, error)
		UpdateAdminEditParamNoCache(ctx context.Context, id string, memo *string, status string) error
		FindUserDetailByIdNoCache(ctx context.Context, id string) (*UserDetail, error)
		InsertUserNoCache(ctx context.Context, id, lineId, firebaseUID, phoneNumber, name string, typeId, genderId uint64, birthday time.Time) error
		UpdateUserEditParamNoCache(ctx context.Context, id string, genderId uint64, birthday time.Time, phoneNumber string, name string) error
		UpdateLastAccessAt(ctx context.Context, id string) error
		UnregisterNoCache(ctx context.Context, id string) error
	}

	customUsersModel struct {
		*defaultUsersModel
	}

	UserWithLineInfo struct {
		Id                  string         `db:"id"`
		FirebaseUid         string         `db:"firebase_uid"`
		Name                string         `db:"user_name"`
		TypeName            string         `db:"type_name"`
		Status              string         `db:"status"`
		CsStatus            string         `db:"cs_status"`
		Memo                sql.NullString `db:"memo"`
		CreatedAt           time.Time      `db:"created_at"`
		LineId              sql.NullString `db:"line_id"`
		IconUrl             sql.NullString `db:"icon_url"`
		LineIdUnregistered  sql.NullString `db:"line_id_unregistered"`
		IconUrlUnregistered sql.NullString `db:"icon_url_unregistered"`
		EmailUnregistered   sql.NullString `db:"email_unregistered"`
	}

	UserDetailForAdmin struct {
		Id                   string         `db:"id"`
		FirebaseUid          string         `db:"firebase_uid"`
		TypeName             string         `db:"type_name"`
		Gender               string         `db:"gender"`
		Name                 string         `db:"name"`
		Nickname             string         `db:"nickname"`
		PhoneNumber          string         `db:"phone_number"`
		University           string         `db:"university"`
		Faculty              string         `db:"faculty"`
		Department           string         `db:"department"`
		GraduationYear       uint64         `db:"graduation_year"`
		Birthday             time.Time      `db:"birthday"`
		PostalCode           string         `db:"postal_code"`
		Address              string         `db:"address"`
		CreatedAt            time.Time      `db:"created_at"`
		ServiceTrigger       string         `db:"service_trigger"`
		ServiceTriggerDetail sql.NullString `db:"service_trigger_detail"`
		Introducer           sql.NullString `db:"introducer"`
		Memo                 sql.NullString `db:"memo"`
		Status               string         `db:"status"`
		CsStatus             string         `db:"cs_status"`
		DesiredAnnualIncome  string         `db:"desired_annual_income"`
		LineId               sql.NullString `db:"line_id"`
		IconUrl              sql.NullString `db:"icon_url"`
		LineIdUnregistered   sql.NullString `db:"line_id_unregistered"`
		IconUrlUnregistered  sql.NullString `db:"icon_url_unregistered"`
		EmailUnregistered    sql.NullString `db:"email_unregistered"`
	}

	UserDetail struct {
		FirebaseUid              string         `db:"firebase_uid"`
		Name                     string         `db:"name"`
		Nickname                 string         `db:"nickname"`
		GenderName               string         `db:"gender_name"`
		University               string         `db:"university"`
		Faculty                  string         `db:"faculty"`
		Department               string         `db:"department"`
		GraduationYear           uint64         `db:"graduation_year"`
		Birthday                 time.Time      `db:"birthday"`
		PhoneNumber              string         `db:"phone_number"`
		PostalCode               string         `db:"postal_code"`
		Address                  string         `db:"address"`
		ServiceTriggerId         uint64         `db:"service_trigger_id"`
		ServiceTriggerName       string         `db:"service_trigger_name"`
		ServiceTriggerDetail     sql.NullString `db:"service_trigger_detail"`
		DesiredAnnualIncomeId    uint64         `db:"desired_annual_income_id"`
		DesiredAnnualIncomeValue string         `db:"desired_annual_income_value"`
		Introducer               sql.NullString `db:"introducer"`
		LineId                   string         `db:"line_id"`
		LineName                 string         `db:"line_name"`
		IconUrl                  sql.NullString `db:"icon_url"`
		StatusMessage            sql.NullString `db:"status_message"`
	}
)

const (
	UserStatusRegistered   = "registered"
	UserStatusBanned       = "banned"
	UserStatusUnregistered = "unregistered"
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c),
	}
}

func (m *customUsersModel) FindUserDetailForAdminByIdNoCache(id string) (user *UserDetailForAdmin, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserDetailForAdminByIdNoCache error")
	}()

	var resp UserDetailForAdmin

	query := `
SELECT
	u.id,
	u.firebase_uid,
	g.name as gender,
	ut.name AS type_name,
	u.name,
	gup.nickname,
	gup.university,
	gup.faculty,
	gup.department,
	gup.graduation_year,
	u.birthday,
	u.phone_number,
	ua.postal_code,
	ua.address,
	st.name AS service_trigger,
	gup.service_trigger_detail,
	gup.introducer,
	u.created_at,
	u.memo,
	u.status,
	ucs.status AS cs_status,
	dai.value AS desired_annual_income,
	la.line_id,
	la.icon_url,
	la_unregistered.line_id AS line_id_unregistered,
	la_unregistered.icon_url AS icon_url_unregistered,
	uup.email AS email_unregistered
FROM users AS u
	INNER JOIN user_addresses AS ua
	ON ua.user_id = u.id
	INNER JOIN general_user_profiles AS gup
	ON gup.user_id = u.id
	INNER JOIN service_triggers AS st
	ON st.id = gup.service_trigger_id
	INNER JOIN desired_annual_incomes AS dai
	ON dai.id = gup.desired_annual_income_id
	INNER JOIN user_types AS ut
	ON ut.id = u.type_id
	INNER JOIN genders AS g
	ON g.id = u.gender_id
	INNER JOIN user_cs_statuses AS ucs
	ON ucs.user_id = u.id
	LEFT JOIN line_accounts AS la
	ON la.line_id = u.line_id
	LEFT JOIN unregistered_user_profiles AS uup
	ON uup.user_id = u.id
	LEFT JOIN line_accounts AS la_unregistered
	ON la_unregistered.line_id = uup.line_id
WHERE u.id = ?
`
	err = m.QueryRowNoCache(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindByIdNoCache(userId string) (user *Users, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByIdNoCache error")
	}()

	var resp Users
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", usersRows, m.table)
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

func (m *customUsersModel) FindOneByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (user *Users, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneByFirebaseUidNoCache error")
	}()

	var resp Users
	query := fmt.Sprintf("select %s from %s where `firebase_uid` = ? limit 1", usersRows, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, firebaseUid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// 退会済ユーザーは除外
func (m *customUsersModel) FindByLineIdExceptUnregisteredNoCache(ctx context.Context, lineId string) (user *Users, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByLineIdExceptUnregisteredNoCache error")
	}()

	var resp Users
	query := fmt.Sprintf("select %s from %s where `line_id` = ? and `status` != ? limit 1", usersRows, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, lineId, UserStatusUnregistered)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindUserWithLineInfoListAll(ctx context.Context) (users []UserWithLineInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserWithLineInfoListAll error")
	}()

	var resp []UserWithLineInfo

	query := `
SELECT
	u.id,
	u.firebase_uid,
	u.name AS user_name,
	ut.name AS type_name,
	u.status,
	ucs.status AS cs_status,
	u.memo,
	u.created_at,
	la.line_id,
	la.icon_url,
	la_unregistered.line_id AS line_id_unregistered,
	la_unregistered.icon_url AS icon_url_unregistered,
	uup.email AS email_unregistered
FROM users AS u
	INNER JOIN user_types AS ut
	ON ut.id = u.type_id
	INNER JOIN user_cs_statuses AS ucs
	ON ucs.user_id = u.id
	LEFT JOIN line_accounts AS la
	ON la.line_id = u.line_id
	LEFT JOIN unregistered_user_profiles AS uup
	ON uup.user_id = u.id
	LEFT JOIN line_accounts AS la_unregistered
	ON la_unregistered.line_id = uup.line_id
ORDER BY u.id DESC
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindUserWithLineInfoListByCsStatus(ctx context.Context, csStatus string) (users []UserWithLineInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserWithLineInfoListByCsStatus error")
	}()

	var resp []UserWithLineInfo

	query := `
SELECT
	u.id,
	u.firebase_uid,
	u.name AS user_name,
	ut.name AS type_name,
	u.status,
	ucs.status AS cs_status,
	u.memo,
	u.created_at,
	la.line_id,
	la.icon_url,
	la_unregistered.line_id AS line_id_unregistered,
	la_unregistered.icon_url AS icon_url_unregistered,
	uup.email AS email_unregistered
FROM users AS u
	INNER JOIN user_types AS ut
	ON ut.id = u.type_id
	INNER JOIN user_cs_statuses AS ucs
	ON ucs.user_id = u.id
	LEFT JOIN line_accounts AS la
	ON la.line_id = u.line_id
	LEFT JOIN unregistered_user_profiles AS uup
	ON uup.user_id = u.id
	LEFT JOIN line_accounts AS la_unregistered
	ON la_unregistered.line_id = uup.line_id
WHERE ucs.status = ?
ORDER BY u.id DESC
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, csStatus)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindSpecificRangeUserWithLineInfoList(ctx context.Context, limit uint64, offset uint64) (users []UserWithLineInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeUserWithLineInfoList error")
	}()

	var resp []UserWithLineInfo

	query := `
SELECT
	u.id,
	u.firebase_uid,
	u.name AS user_name,
	ut.name AS type_name,
	u.status,
	ucs.status AS cs_status,
	u.memo,
	u.created_at,
	la.line_id,
	la.icon_url,
	la_unregistered.line_id AS line_id_unregistered,
	la_unregistered.icon_url AS icon_url_unregistered,
	uup.email AS email_unregistered
FROM users AS u
	INNER JOIN user_types AS ut
	ON ut.id = u.type_id
	INNER JOIN user_cs_statuses AS ucs
	ON ucs.user_id = u.id
	LEFT JOIN line_accounts AS la
	ON la.line_id = u.line_id
	LEFT JOIN unregistered_user_profiles AS uup
	ON uup.user_id = u.id
	LEFT JOIN line_accounts AS la_unregistered
	ON la_unregistered.line_id = uup.line_id
ORDER BY u.id DESC
LIMIT ? OFFSET ?
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, limit, offset)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindSpecificRangeUserWithLineInfoListByCsStatus(ctx context.Context, csStatus string, limit uint64, offset uint64) (users []UserWithLineInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeUserWithLineInfoListByCsStatus error")
	}()

	var resp []UserWithLineInfo

	query := `
SELECT
	u.id,
	u.firebase_uid,
	u.name AS user_name,
	ut.name AS type_name,
	u.status,
	ucs.status AS cs_status,
	u.memo,
	u.created_at,
	la.line_id,
	la.icon_url,
	la_unregistered.line_id AS line_id_unregistered,
	la_unregistered.icon_url AS icon_url_unregistered,
	uup.email AS email_unregistered
FROM users AS u
	INNER JOIN user_types AS ut
	ON ut.id = u.type_id
	INNER JOIN user_cs_statuses AS ucs
	ON ucs.user_id = u.id
	LEFT JOIN line_accounts AS la
	ON la.line_id = u.line_id
	LEFT JOIN unregistered_user_profiles AS uup
	ON uup.user_id = u.id
	LEFT JOIN line_accounts AS la_unregistered
	ON la_unregistered.line_id = uup.line_id
WHERE ucs.status = ?
ORDER BY u.id DESC
LIMIT ? OFFSET ?
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, csStatus, limit, offset)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindUserDetailByIdNoCache(ctx context.Context, id string) (user *UserDetail, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserDetailByIdNoCache error")
	}()

	var resp UserDetail

	query := `
SELECT
	u.firebase_uid,
	u.name,
	gup.nickname,
	gup.university,
	gup.faculty,
	gup.department,
	gup.graduation_year,
	g.name as gender_name,
	ua.postal_code,
	ua.address,
	dai.value as desired_annual_income,
	st.id as service_trigger_id,
	st.name as service_trigger_name,
	gup.service_trigger_detail,
	gup.introducer,
	dai.id as desired_annual_income_id,
	dai.value as desired_annual_income_value,
	u.birthday,
	u.phone_number,
	u.line_id,
    la.name AS line_name,
    la.icon_url,
    la.status_message
FROM users AS u
	INNER JOIN user_addresses AS ua
	ON ua.user_id = u.id
	INNER JOIN general_user_profiles AS gup
	ON gup.user_id = u.id
	INNER JOIN genders AS g
	ON g.id = u.gender_id
	INNER JOIN desired_annual_incomes AS dai
	ON dai.id = gup.desired_annual_income_id
    INNER JOIN line_accounts AS la
    ON la.line_id = u.line_id
	INNER JOIN service_triggers AS st
	ON st.id = gup.service_trigger_id
WHERE u.id = ?
`
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, id)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) FindCount(ctx context.Context) (count uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCount error")
	}()

	query := fmt.Sprintf(
		"select COUNT(`id`) AS `count` from %s",
		m.table,
	)

	err = m.QueryRowNoCacheCtx(ctx, &count, query)

	switch err {
	case nil:
		return count, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *customUsersModel) FindCountByCsStatus(ctx context.Context, csStatus string) (count uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCountByCsStatus error")
	}()

	query := `
	SELECT COUNT(*)
	FROM users AS u
		INNER JOIN user_cs_statuses AS ucs
		ON ucs.user_id = u.id
	WHERE ucs.status = ?
	`

	err = m.QueryRowNoCacheCtx(ctx, &count, query, csStatus)

	switch err {
	case nil:
		return count, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *customUsersModel) UpdateAdminEditParamNoCache(ctx context.Context, id string, memo *string, status string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateAdminEditParamNoCache error")
	}()

	query := `
UPDATE
users
SET memo = ?, status = ?
WHERE id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		memo,
		status,
		id,
	)

	return err
}

func (m *customUsersModel) InsertUserNoCache(
	ctx context.Context,
	id,
	lineId,
	firebaseUID,
	phoneNumber,
	name string,
	typeId,
	genderId uint64,
	birthday time.Time,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "InsertUserNoCache error")
	}()

	query := `
INSERT INTO users(
	id,
	line_id,
	firebase_uid,
	type_id,
	gender_id,
	birthday,
	phone_number,
	name
) VALUES(?, ?, ?, ?, ?, ?, ?, ?)
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		id,
		lineId,
		firebaseUID,
		typeId,
		genderId,
		birthday,
		phoneNumber,
		name,
	)
	return err
}

func (m *customUsersModel) UpdateUserEditParamNoCache(
	ctx context.Context,
	id string,
	genderId uint64,
	birthday time.Time,
	phoneNumber string,
	name string,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateUserEditParamNoCache error")
	}()

	query := `
UPDATE
users
SET gender_id = ?, birthday = ?, phone_number = ?, name = ?
WHERE id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		genderId,
		birthday,
		phoneNumber,
		name,
		id,
	)
	return err
}

func (m *customUsersModel) UpdateLastAccessAt(ctx context.Context, id string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateLastAccessAt error")
	}()

	query := fmt.Sprintf("UPDATE %s SET `last_access_at` = CURRENT_TIMESTAMP() WHERE `id` = ?", m.table)
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		id,
	)
	return err
}

func (m *customUsersModel) UnregisterNoCache(ctx context.Context, id string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UnregisterNoCache error")
	}()

	query := fmt.Sprintf("UPDATE %s SET `status` = ?, `line_id` = NULL WHERE `id` = ?", m.table)
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		UserStatusUnregistered,
		id,
	)
	return err
}
