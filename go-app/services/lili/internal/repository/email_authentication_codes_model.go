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

var _ EmailAuthenticationCodesModel = (*customEmailAuthenticationCodesModel)(nil)

type (
	// EmailAuthenticationCodesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailAuthenticationCodesModel.
	EmailAuthenticationCodesModel interface {
		emailAuthenticationCodesModel
		FindByInactiveUserIdNoCache(ctx context.Context, inactiveUserId string) (*EmailAuthenticationCodes, error)
		InsertNoCache(ctx context.Context, data *EmailAuthenticationCodes) (sql.Result, error)
		UpdateAttemptCountNoCache(ctx context.Context, id string, attemptCount uint64) (sql.Result, error)
		UpdateCodeHashNoCache(ctx context.Context, id, codeHash string) (sql.Result, error)
	}

	customEmailAuthenticationCodesModel struct {
		*defaultEmailAuthenticationCodesModel
	}
)

// NewEmailAuthenticationCodesModel returns a model for the database table.
func NewEmailAuthenticationCodesModel(conn sqlx.SqlConn, c cache.CacheConf) EmailAuthenticationCodesModel {
	return &customEmailAuthenticationCodesModel{
		defaultEmailAuthenticationCodesModel: newEmailAuthenticationCodesModel(conn, c),
	}
}

func (m *customEmailAuthenticationCodesModel) FindByInactiveUserIdNoCache(ctx context.Context, inactiveUserId string) (c *EmailAuthenticationCodes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByInactiveUserIdNoCache error")
	}()

	var resp EmailAuthenticationCodes
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `inactive_user_id` = ? ORDER BY `id` DESC LIMIT 1", emailAuthenticationCodesRows, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, inactiveUserId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEmailAuthenticationCodesModel) InsertNoCache(ctx context.Context, data *EmailAuthenticationCodes) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	rows := strings.Join(stringx.Remove(emailAuthenticationCodesFieldNames, "`created_at`", "`updated_at`"), ",")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?, ?, ?, ?)", m.table, rows)
	ret, err := m.ExecNoCacheCtx(
		ctx,
		query,
		data.Id,
		data.InactiveUserId,
		data.CodeHash,
		data.AttemptCount,
	)

	return ret, err
}

func (m *customEmailAuthenticationCodesModel) UpdateAttemptCountNoCache(ctx context.Context, id string, attemptCount uint64) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateAttemptCountNoCache error")
	}()

	query := fmt.Sprintf("UPDATE %s SET `attempt_count` = ? WHERE `id` = ?", m.table)
	ret, err := m.ExecNoCacheCtx(
		ctx,
		query,
		attemptCount,
		id,
	)

	return ret, err
}

func (m *customEmailAuthenticationCodesModel) UpdateCodeHashNoCache(ctx context.Context, id, codeHash string) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateCodeHashNoCache error")
	}()

	query := fmt.Sprintf("UPDATE %s SET `code_hash` = ? WHERE `id` = ?", m.table)
	ret, err := m.ExecNoCacheCtx(
		ctx,
		query,
		codeHash,
		id,
	)

	return ret, err
}
