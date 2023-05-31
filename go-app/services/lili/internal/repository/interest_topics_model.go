package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InterestTopicsModel = (*customInterestTopicsModel)(nil)

type (
	// InterestTopicsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInterestTopicsModel.
	InterestTopicsModel interface {
		interestTopicsModel
		FindAll(ctx context.Context) ([]InterestTopics, error)
	}

	customInterestTopicsModel struct {
		*defaultInterestTopicsModel
	}
)

// NewInterestTopicsModel returns a model for the database table.
func NewInterestTopicsModel(conn sqlx.SqlConn, c cache.CacheConf) InterestTopicsModel {
	return &customInterestTopicsModel{
		defaultInterestTopicsModel: newInterestTopicsModel(conn, c),
	}
}

func (m *customInterestTopicsModel) FindAll(ctx context.Context) (topics []InterestTopics, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []InterestTopics
	cacheKey := "cache:interestTopics:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s order by display_order", interestTopicsRows, m.table)
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
