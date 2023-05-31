package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ InactiveUsersModel = (*customInactiveUsersModel)(nil)

type (
	// InactiveUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInactiveUsersModel.
	InactiveUsersModel interface {
		inactiveUsersModel
		FindByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (*InactiveUsers, error)
		FindWithLineInfoByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (*InactiveUsersWithLine, error)
		FindByLineIdNoCache(ctx context.Context, lineId string) (*InactiveUsers, error)
		Upsert(ctx context.Context, data *InactiveUsers) (*InactiveUsers, error)
	}

	customInactiveUsersModel struct {
		*defaultInactiveUsersModel
	}

	InactiveUsersWithLine struct {
		Id            string         `db:"id"`
		TypeId        uint64         `db:"type_id"`
		FirebaseUid   string         `db:"firebase_uid"`
		Name          string         `db:"name"`
		LineId        string         `db:"line_id"`
		LineName      string         `db:"line_name"`
		IconUrl       sql.NullString `db:"icon_url"`
		StatusMessage sql.NullString `db:"status_message"`
	}
)

// NewInactiveUsersModel returns a model for the database table.
func NewInactiveUsersModel(conn sqlx.SqlConn, c cache.CacheConf) InactiveUsersModel {
	return &customInactiveUsersModel{
		defaultInactiveUsersModel: newInactiveUsersModel(conn, c),
	}
}

func (m *customInactiveUsersModel) FindByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (u *InactiveUsers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByFirebaseUIDNoCache error")
	}()

	var resp InactiveUsers
	query := fmt.Sprintf("select %s from %s where `firebase_uid` = ? limit 1", inactiveUsersRows, m.table)
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

func (m *customInactiveUsersModel) FindWithLineInfoByFirebaseUidNoCache(ctx context.Context, firebaseUid string) (user *InactiveUsersWithLine, err error) {
	defer func() {
		err = errors.Wrap(err, "FindWithLineInfoByFirebaseUidNoCache error")
	}()

	var resp InactiveUsersWithLine
	query := fmt.Sprintf(`
SELECT 
  iu.id,
  iu.type_id,
  iu.firebase_uid,
  iu.name,
  iu.line_id,
  la.name AS line_name,
  la.icon_url,
  la.status_message
FROM %s AS iu
  INNER JOIN line_accounts AS la
  ON la.line_id = iu.line_id
WHERE iu.firebase_uid = ? limit 1
`, m.table)

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

func (m *customInactiveUsersModel) Upsert(ctx context.Context, data *InactiveUsers) (user *InactiveUsers, err error) {
	defer func() {
		err = errors.Wrap(err, "Upsert error")
	}()

	columns := strings.Join(stringx.Remove(inactiveUsersFieldNames, "`created_at`", "`updated_at`"), ",")
	upsertSql := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE
			type_id = VALUES(type_id),
			name = VALUES(name),
			firebase_uid = VALUES(firebase_uid)`, // 退会後の再登録ではfirebaseユーザが新規作成されるのでfirebase_uidの更新が必要
		m.table, columns)

	if _, err = m.ExecNoCacheCtx(
		ctx,
		upsertSql,
		data.Id,
		data.LineId,
		data.TypeId,
		data.FirebaseUid,
		data.Name,
	); err != nil {
		return nil, err
	}

	user, err = m.FindByFirebaseUidNoCache(ctx, data.FirebaseUid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *customInactiveUsersModel) FindByLineIdNoCache(ctx context.Context, lineId string) (user *InactiveUsers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByLineIdNoCache error")
	}()

	var resp InactiveUsers
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `line_id` = ? LIMIT 1", inactiveUsersRows, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, lineId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
