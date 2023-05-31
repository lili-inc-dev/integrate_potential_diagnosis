package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ LineAccountsModel = (*customLineAccountsModel)(nil)

type (
	// LineAccountsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLineAccountsModel.
	LineAccountsModel interface {
		lineAccountsModel
		Upsert(ctx context.Context, lineId, name, iconURL, statusMsg string) (sql.Result, error)
	}

	customLineAccountsModel struct {
		*defaultLineAccountsModel
	}
)

// NewLineAccountsModel returns a model for the database table.
func NewLineAccountsModel(conn sqlx.SqlConn, c cache.CacheConf) LineAccountsModel {
	return &customLineAccountsModel{
		defaultLineAccountsModel: newLineAccountsModel(conn, c),
	}
}

func (m *customLineAccountsModel) Upsert(ctx context.Context, lineId, name, iconURL, statusMsg string) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "Upsert error")
	}()

	rows := strings.Join(stringx.Remove(lineAccountsFieldNames, "`created_at`", "`updated_at`"), ",")
	upsertSql := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE
			name = VALUES(name), 
			icon_url = VALUES(icon_url),
			status_message = VALUES(status_message)`,
		m.table, rows)

	ret, err := m.ExecNoCacheCtx(
		ctx,
		upsertSql,
		lineId, name, iconURL, statusMsg,
	)

	return ret, err
}
