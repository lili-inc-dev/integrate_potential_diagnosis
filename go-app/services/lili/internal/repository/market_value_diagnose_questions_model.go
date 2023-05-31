package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketValueDiagnoseQuestionsModel = (*customMarketValueDiagnoseQuestionsModel)(nil)

type (
	// MarketValueDiagnoseQuestionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketValueDiagnoseQuestionsModel.
	MarketValueDiagnoseQuestionsModel interface {
		marketValueDiagnoseQuestionsModel
		FindCount(ctx context.Context) (*uint64, error)
		FindListByDisplayOrder(ctx context.Context, displayOrder uint64) ([]MarketValueDiagnoseQuestions, error)
	}

	customMarketValueDiagnoseQuestionsModel struct {
		*defaultMarketValueDiagnoseQuestionsModel
	}
)

// NewMarketValueDiagnoseQuestionsModel returns a model for the database table.
func NewMarketValueDiagnoseQuestionsModel(conn sqlx.SqlConn, c cache.CacheConf) MarketValueDiagnoseQuestionsModel {
	return &customMarketValueDiagnoseQuestionsModel{
		defaultMarketValueDiagnoseQuestionsModel: newMarketValueDiagnoseQuestionsModel(conn, c),
	}
}

func (m *customMarketValueDiagnoseQuestionsModel) FindCount(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCount error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}

	cacheKey := "cache:marketValueDiagnoseQuestions:count"
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

func (m *customMarketValueDiagnoseQuestionsModel) FindListByDisplayOrder(ctx context.Context, displayOrder uint64) (questions []MarketValueDiagnoseQuestions, err error) {
	defer func() {
		err = errors.Wrap(err, "FindListByDisplayOrder error")
	}()

	var resp []MarketValueDiagnoseQuestions
	cacheKey := "cache:marketValueDiagnoseQuestions:displayOrder:" + strconv.FormatUint(displayOrder, 10) + ":"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf(
			`
SELECT
	mvdq.id,
	mvdq.market_value_id,
	mvdq.index,
	mvdq.content,
	mvdq.created_at,
	mvdq.updated_at
FROM %s AS mvdq
	INNER JOIN market_values AS m
	ON m.id = mvdq.market_value_id
	AND m.display_order = ?
ORDER BY mvdq.index 
`,
			m.table,
		)
		return conn.QueryRows(&resp, query, strconv.FormatUint(displayOrder, 10))
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
