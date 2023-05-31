package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GendersModel = (*customGendersModel)(nil)

type (
	// GendersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGendersModel.
	GendersModel interface {
		gendersModel
		FindAll(ctx context.Context) ([]Genders, error)
	}

	customGendersModel struct {
		*defaultGendersModel
	}
)

// NewGendersModel returns a model for the database table.
func NewGendersModel(conn sqlx.SqlConn, c cache.CacheConf) GendersModel {
	return &customGendersModel{
		defaultGendersModel: newGendersModel(conn, c),
	}
}

func (m *customGendersModel) FindAll(ctx context.Context) (genders []Genders, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []Genders
	cacheKey := "cache:genders:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s order by id", gendersRows, m.table)
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
