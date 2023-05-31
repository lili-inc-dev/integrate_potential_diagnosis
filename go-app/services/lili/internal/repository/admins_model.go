package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ AdminsModel = (*customAdminsModel)(nil)

type (
	// AdminsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminsModel.
	AdminsModel interface {
		adminsModel
		FindOneNoCache(ctx context.Context, id uint64) (*Admins, error)
		FindCount(ctx context.Context) (*uint64, error)
		FindOneByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (*Admins, error)
		InsertNoCache(ctx context.Context, data *Admins) (sql.Result, error)
		FindOneWithRoleNoCache(ctx context.Context, id uint64) (*AdminWithRoleInfo, error)
		FindSpecificRangeListWithRoleNoCache(ctx context.Context, limit uint64, offset uint64) ([]AdminWithRoleInfo, error)
		UpdateLastAccessAt(ctx context.Context, id uint64) error
	}

	customAdminsModel struct {
		*defaultAdminsModel
	}

	AdminWithRoleInfo struct {
		Id           uint64         `db:"id"`
		FirebaseUid  string         `db:"firebase_uid"`
		RoleId       string         `db:"role_id"`
		RoleName     string         `db:"role_name"`
		Name         string         `db:"name"`
		LastAccessAt time.Time      `db:"last_access_at"`
		Affiliation  sql.NullString `db:"affiliation"`
		IsDisabled   bool           `db:"is_disabled"`
	}
)

// NewAdminsModel returns a model for the database table.
func NewAdminsModel(conn sqlx.SqlConn, c cache.CacheConf) AdminsModel {
	return &customAdminsModel{
		defaultAdminsModel: newAdminsModel(conn, c),
	}
}

func (m *customAdminsModel) FindOneNoCache(ctx context.Context, id uint64) (admin *Admins, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneNoCache error")
	}()

	var resp Admins
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", adminsRows, m.table)
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

func (m *customAdminsModel) FindCount(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCount error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}

	query := fmt.Sprintf(
		"select COUNT(`id`) AS `count` from %s",
		m.table,
	)

	err = m.QueryRowNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return &resp.Count, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customAdminsModel) FindOneByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (admin *Admins, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneByFirebaseUidNoCache error")
	}()

	var resp Admins
	query := fmt.Sprintf("select %s from %s where `firebase_uid` = ? limit 1", adminsRows, m.table)
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

func (m *customAdminsModel) FindOneWithRoleNoCache(ctx context.Context, id uint64) (admin *AdminWithRoleInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneWithRoleNoCache error")
	}()

	var resp AdminWithRoleInfo

	query := `
SELECT
	a.id,
	a.firebase_uid,
	a.name,
	ar.id AS role_id,
	ar.name AS role_name,
	a.last_access_at,
	a.affiliation,
	a.is_disabled
FROM admins AS a
	INNER JOIN admin_roles AS ar
	ON ar.id = a.role_id
WHERE a.id = ?
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

func (m *customAdminsModel) InsertNoCache(ctx context.Context, data *Admins) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	rows := strings.Join(stringx.Remove(adminsFieldNames, "`id`", "`last_access_at`", "`created_at`", "`updated_at`"), ",")
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, rows)
	ret, err := m.ExecNoCacheCtx(
		ctx,
		query,
		data.RoleId,
		data.FirebaseUid,
		data.Name,
		data.Affiliation,
		data.IsDisabled,
	)

	return ret, err
}

func (m *customAdminsModel) FindSpecificRangeListWithRoleNoCache(ctx context.Context, limit uint64, offset uint64) (admins []AdminWithRoleInfo, err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	var resp []AdminWithRoleInfo

	query := `
SELECT
	a.id,
	a.firebase_uid,
	a.name,
	ar.id AS role_id,
	ar.name AS role_name,
	a.last_access_at,
	a.affiliation,
	a.is_disabled
FROM admins AS a
	INNER JOIN admin_roles AS ar
	ON ar.id = a.role_id
ORDER BY a.id DESC
LIMIT ?
OFFSET ?
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

func (m *customAdminsModel) Update(ctx context.Context, data *Admins) (err error) {
	defer func() {
		err = errors.Wrap(err, "Update error")
	}()

	rowsWithPlaceHolder := strings.Join(stringx.Remove(adminsFieldNames, "`id`", "`firebase_uid`", "`email`", "`last_access_at`", "`created_at`"), "=?,") + "=?"
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `id` = ?", m.table, rowsWithPlaceHolder)
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		data.RoleId,
		data.Name,
		data.Affiliation,
		data.IsDisabled,
		time.Now(),
		data.Id,
	)
	return err
}

func (m *customAdminsModel) UpdateLastAccessAt(ctx context.Context, id uint64) (err error) {
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
