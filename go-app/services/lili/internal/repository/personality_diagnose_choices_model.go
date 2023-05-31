package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PersonalityDiagnoseChoicesModel = (*customPersonalityDiagnoseChoicesModel)(nil)

type (
	// PersonalityDiagnoseChoicesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPersonalityDiagnoseChoicesModel.
	PersonalityDiagnoseChoicesModel interface {
		personalityDiagnoseChoicesModel
		FindAll(ctx context.Context) ([]PersonalityDiagnoseChoices, error)
	}

	customPersonalityDiagnoseChoicesModel struct {
		*defaultPersonalityDiagnoseChoicesModel
	}
)

// NewPersonalityDiagnoseChoicesModel returns a model for the database table.
func NewPersonalityDiagnoseChoicesModel(conn sqlx.SqlConn, c cache.CacheConf) PersonalityDiagnoseChoicesModel {
	return &customPersonalityDiagnoseChoicesModel{
		defaultPersonalityDiagnoseChoicesModel: newPersonalityDiagnoseChoicesModel(conn, c),
	}
}

func (m *customPersonalityDiagnoseChoicesModel) FindAll(ctx context.Context) (choices []PersonalityDiagnoseChoices, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAll error")
	}()

	var resp []PersonalityDiagnoseChoices
	cacheKey := "cache:personalityDiagnoseQuestions:"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s", personalityDiagnoseChoicesRows, m.table)
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
