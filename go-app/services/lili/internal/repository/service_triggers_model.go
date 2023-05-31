package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ServiceTriggersModel = (*customServiceTriggersModel)(nil)

type (
	// ServiceTriggersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServiceTriggersModel.
	ServiceTriggersModel interface {
		serviceTriggersModel
		FindAll(ctx context.Context) ([]ServiceTriggers, error)
	}

	customServiceTriggersModel struct {
		*defaultServiceTriggersModel
	}
)

// NewServiceTriggersModel returns a model for the database table.
func NewServiceTriggersModel(conn sqlx.SqlConn, c cache.CacheConf) ServiceTriggersModel {
	return &customServiceTriggersModel{
		defaultServiceTriggersModel: newServiceTriggersModel(conn, c),
	}
}

func (m *customServiceTriggersModel) FindAll(ctx context.Context) (triggers []ServiceTriggers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []ServiceTriggers
	cacheKey := "cache:serviceTriggers:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s order by id", serviceTriggersRows, m.table)
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
