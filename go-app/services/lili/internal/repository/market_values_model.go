package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketValuesModel = (*customMarketValuesModel)(nil)

type (
	// MarketValuesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketValuesModel.
	MarketValuesModel interface {
		marketValuesModel
		FindCount(ctx context.Context) (*uint64, error)
	}

	customMarketValuesModel struct {
		*defaultMarketValuesModel
	}
)

// NewMarketValuesModel returns a model for the database table.
func NewMarketValuesModel(conn sqlx.SqlConn, c cache.CacheConf) MarketValuesModel {
	return &customMarketValuesModel{
		defaultMarketValuesModel: newMarketValuesModel(conn, c),
	}
}

func (m *customMarketValuesModel) FindCount(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCount")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}

	cacheKey := "cache:marketValues:count"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf(
			"select COUNT(`id`) AS `count` from %s",
			m.table,
		)
		return conn.QueryRow(&resp, query)
	})

	switch err {
	case nil:
		return &resp.Count, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
