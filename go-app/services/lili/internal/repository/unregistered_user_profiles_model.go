package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ UnregisteredUserProfilesModel = (*customUnregisteredUserProfilesModel)(nil)

type (
	// UnregisteredUserProfilesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUnregisteredUserProfilesModel.
	UnregisteredUserProfilesModel interface {
		unregisteredUserProfilesModel
		InsertNoCache(ctx context.Context, id, userId, lineId, email string, reasonId uint64) error
	}

	customUnregisteredUserProfilesModel struct {
		*defaultUnregisteredUserProfilesModel
	}
)

// NewUnregisteredUserProfilesModel returns a model for the database table.
func NewUnregisteredUserProfilesModel(conn sqlx.SqlConn, c cache.CacheConf) UnregisteredUserProfilesModel {
	return &customUnregisteredUserProfilesModel{
		defaultUnregisteredUserProfilesModel: newUnregisteredUserProfilesModel(conn, c),
	}
}

func (m *customUnregisteredUserProfilesModel) InsertNoCache(ctx context.Context, id, userId, lineId, email string, reasonId uint64) (err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	rows := strings.Join(stringx.Remove(unregisteredUserProfilesFieldNames, "`created_at`", "`updated_at`"), ",")
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?, ?, ?, ?, ?)`, m.table, rows)
	_, err = m.ExecNoCacheCtx(ctx, query, id, userId, reasonId, lineId, email)
	return err
}
