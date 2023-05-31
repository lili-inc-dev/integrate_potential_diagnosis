package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DesiredAnnualIncomesModel = (*customDesiredAnnualIncomesModel)(nil)

type (
	// DesiredAnnualIncomesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDesiredAnnualIncomesModel.
	DesiredAnnualIncomesModel interface {
		desiredAnnualIncomesModel
		FindAll(ctx context.Context) ([]DesiredAnnualIncomes, error)
	}

	customDesiredAnnualIncomesModel struct {
		*defaultDesiredAnnualIncomesModel
	}
)

// NewDesiredAnnualIncomesModel returns a model for the database table.
func NewDesiredAnnualIncomesModel(conn sqlx.SqlConn, c cache.CacheConf) DesiredAnnualIncomesModel {
	return &customDesiredAnnualIncomesModel{
		defaultDesiredAnnualIncomesModel: newDesiredAnnualIncomesModel(conn, c),
	}
}

func (m *customDesiredAnnualIncomesModel) FindAll(ctx context.Context) (incomes []DesiredAnnualIncomes, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []DesiredAnnualIncomes

	cacheKey := "cache:desiredAnnualIncomes:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s order by id", desiredAnnualIncomesRows, m.table)
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
