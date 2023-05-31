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

var _ PersonalityDiagnoseQuestionsModel = (*customPersonalityDiagnoseQuestionsModel)(nil)

type (
	// PersonalityDiagnoseQuestionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPersonalityDiagnoseQuestionsModel.
	PersonalityDiagnoseQuestionsModel interface {
		personalityDiagnoseQuestionsModel
		FindCount(ctx context.Context) (*uint64, error)
		FindSpecificRangeList(ctx context.Context, limit uint64, offset uint64) ([]PersonalityDiagnoseQuestions, error)
		FindPerfectPointList(ctx context.Context) ([]PersonalityPerfectPoint, error)
	}

	customPersonalityDiagnoseQuestionsModel struct {
		*defaultPersonalityDiagnoseQuestionsModel
	}

	PersonalityPerfectPoint struct {
		PersonalityName string `db:"personality_name"`
		PerfectPoint    uint64 `db:"perfect_point"`
	}
)

// NewPersonalityDiagnoseQuestionsModel returns a model for the database table.
func NewPersonalityDiagnoseQuestionsModel(conn sqlx.SqlConn, c cache.CacheConf) PersonalityDiagnoseQuestionsModel {
	return &customPersonalityDiagnoseQuestionsModel{
		defaultPersonalityDiagnoseQuestionsModel: newPersonalityDiagnoseQuestionsModel(conn, c),
	}
}

func (m *customPersonalityDiagnoseQuestionsModel) FindCount(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCount error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}

	cacheKey := "cache:personalityDiagnoseQuestions:count"
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

func (m *customPersonalityDiagnoseQuestionsModel) FindSpecificRangeList(ctx context.Context, limit uint64, offset uint64) (questions []PersonalityDiagnoseQuestions, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeList error")
	}()

	var resp []PersonalityDiagnoseQuestions
	cacheKey := "cache:personalityDiagnoseQuestions:offset:" + strconv.FormatUint(offset, 10) + ":limit" + strconv.FormatUint(limit, 10) + ":"
	err = m.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf(
			"select %s from %s order by `index` limit %s offset %s",
			personalityDiagnoseQuestionsRows,
			m.table,
			strconv.FormatUint(limit, 10),
			strconv.FormatUint(offset, 10),
		)
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

func (m *customPersonalityDiagnoseQuestionsModel) FindPerfectPointList(ctx context.Context) (result []PersonalityPerfectPoint, err error) {
	defer func() {
		err = errors.Wrap(err, "FindPerfectPointList error")
	}()

	var resp []PersonalityPerfectPoint

	query := `
SELECT
	COUNT(*) * (
		SELECT weight FROM personality_diagnose_choices ORDER BY weight DESC LIMIT 1
	) AS perfect_point,
	p.name AS personality_name
FROM personality_diagnose_questions AS pdq
	INNER JOIN personalities AS p
	ON pdq.personality_id = p.id
GROUP BY p.id
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
