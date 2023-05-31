package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminRolesModel = (*customAdminRolesModel)(nil)

type (
	// AdminRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminRolesModel.
	AdminRolesModel interface {
		adminRolesModel
		FindAll(ctx context.Context) ([]AdminRoles, error)
		FindAdminRoleWithGrantedCountListAll() ([]AdminRoleWithGrantedCount, error)
	}

	customAdminRolesModel struct {
		*defaultAdminRolesModel
	}

	AdminRoleWithGrantedCount struct {
		Id           uint64         `db:"id"`
		Name         string         `db:"name"`
		Description  sql.NullString `db:"description"`
		GrantedCount uint64         `db:"granted_count"`
	}
)

// NewAdminRolesModel returns a model for the database table.
func NewAdminRolesModel(conn sqlx.SqlConn, c cache.CacheConf) AdminRolesModel {
	return &customAdminRolesModel{
		defaultAdminRolesModel: newAdminRolesModel(conn, c),
	}
}

const findAllCacheKey = "cache:adminRoles"

func (m *customAdminRolesModel) FindAll(ctx context.Context) (resp []AdminRoles, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	cacheKey := findAllCacheKey
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s", adminRolesRows, m.table)
		return conn.QueryRows(&resp, query)
	})
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customAdminRolesModel) FindAdminRoleWithGrantedCountListAll() (resp []AdminRoleWithGrantedCount, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	query := `
SELECT
	ar.id,
	ar.name,
	ar.description,
	COUNT(a.id) AS granted_count
FROM admin_roles AS ar
	LEFT JOIN admins AS a
	ON a.role_id = ar.id
GROUP BY ar.id
`

	err = m.QueryRowsNoCache(&resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
