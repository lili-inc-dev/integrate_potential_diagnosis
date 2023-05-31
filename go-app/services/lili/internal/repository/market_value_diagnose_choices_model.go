package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketValueDiagnoseChoicesModel = (*customMarketValueDiagnoseChoicesModel)(nil)

type (
	// MarketValueDiagnoseChoicesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketValueDiagnoseChoicesModel.
	MarketValueDiagnoseChoicesModel interface {
		marketValueDiagnoseChoicesModel
		FindAll(ctx context.Context) ([]MarketValueDiagnoseChoices, error)
	}

	customMarketValueDiagnoseChoicesModel struct {
		*defaultMarketValueDiagnoseChoicesModel
	}
)

// NewMarketValueDiagnoseChoicesModel returns a model for the database table.
func NewMarketValueDiagnoseChoicesModel(conn sqlx.SqlConn, c cache.CacheConf) MarketValueDiagnoseChoicesModel {
	return &customMarketValueDiagnoseChoicesModel{
		defaultMarketValueDiagnoseChoicesModel: newMarketValueDiagnoseChoicesModel(conn, c),
	}
}

func (m *customMarketValueDiagnoseChoicesModel) FindAll(ctx context.Context) (choices []MarketValueDiagnoseChoices, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []MarketValueDiagnoseChoices
	cacheKey := "cache:marketValueDiagnoseChoices:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s", marketValueDiagnoseChoicesRows, m.table)
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
