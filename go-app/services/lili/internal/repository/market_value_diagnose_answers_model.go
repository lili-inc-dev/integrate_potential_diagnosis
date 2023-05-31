package repository

import (
	"context"
	"fmt"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketValueDiagnoseAnswersModel = (*customMarketValueDiagnoseAnswersModel)(nil)

type (
	// MarketValueDiagnoseAnswersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketValueDiagnoseAnswersModel.
	MarketValueDiagnoseAnswersModel interface {
		marketValueDiagnoseAnswersModel
		BulkInsert(answers []MarketValueDiagnoseAnswers) error
		FindAnswerResultListByAnswerGroupId(ctx context.Context, answerGroupId string) ([]MarketValueDiagnoseAnswerResult, error)
		FindAnswerResultList(ctx context.Context, userId string) ([]MarketValueDiagnoseAnswerResult, error)
		FindLatestAnswerList(ctx context.Context, userId string) (answers []MarketValueDiagnoseAnswers, err error)
	}

	customMarketValueDiagnoseAnswersModel struct {
		*defaultMarketValueDiagnoseAnswersModel
		bulkInserter sqlx.BulkInserter
	}

	MarketValueDiagnoseAnswerResult struct {
		MarketValueName string `db:"market_value_name"`
		TotalPoint      uint64 `db:"total_point"`
	}
)

const (
	// 市場価値診断の満点
	MarketValueDiagnosePerfectPoint uint64 = 10

	marketValueAnswerInsertSql = `
insert into market_value_diagnose_answers 
(id,answer_group_id,question_id,choice_id,user_id)
values (?,?,?,?,?)`

	marketValueMaxWeightByQuestionQuery = `
SELECT
	MAX(mdqw.weight) AS max_weight,
	mdqw.question_id
FROM market_value_diagnose_question_weights AS mdqw
GROUP BY mdqw.question_id
`

	marketValueSumPointByMarketValueQuery = `
SELECT
	SUM(mxwq.max_weight) AS sum_point,
	mvdq.market_value_id
FROM (` + marketValueMaxWeightByQuestionQuery + `) AS mxwq
	INNER JOIN market_value_diagnose_questions AS mvdq
	ON mvdq.id = mxwq.question_id
GROUP BY mvdq.market_value_id
`
)

// NewMarketValueDiagnoseAnswersModel returns a model for the database table.
func NewMarketValueDiagnoseAnswersModel(conn sqlx.SqlConn, c cache.CacheConf) MarketValueDiagnoseAnswersModel {
	bulkInserter, _ := sqlx.NewBulkInserter(conn, marketValueAnswerInsertSql)
	return &customMarketValueDiagnoseAnswersModel{
		defaultMarketValueDiagnoseAnswersModel: newMarketValueDiagnoseAnswersModel(conn, c),
		bulkInserter:                           *bulkInserter,
	}
}

func (m *customMarketValueDiagnoseAnswersModel) BulkInsert(answers []MarketValueDiagnoseAnswers) (err error) {
	defer func() {
		err = errors.Wrap(err, "BulkInsert error")
	}()

	answerGroupId, err := util.GenerateUlid()
	if err != nil {
		return err
	}
	for _, answer := range answers {
		ulid, err := util.GenerateUlid()
		if err != nil {
			return err
		}
		err = m.bulkInserter.Insert(ulid.String(), answerGroupId.String(), answer.QuestionId, answer.ChoiceId, answer.UserId)
		if err != nil {
			return err
		}
	}
	m.bulkInserter.Flush()
	return nil
}

func (m *customMarketValueDiagnoseAnswersModel) FindAnswerResultListByAnswerGroupId(ctx context.Context, answerGroupId string) (results []MarketValueDiagnoseAnswerResult, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAnswerResultListByAnswerGroupId error")
	}()

	var resp []MarketValueDiagnoseAnswerResult

	totalPointByMarketValueQuery := `
SELECT
	SUM(mdqw.weight) AS total_point,
	m.id AS market_value_id,
	m.name AS market_value_name
FROM market_value_diagnose_answers AS mvda
	INNER JOIN market_value_diagnose_questions AS mvdq
	ON mvdq.id = mvda.question_id
	INNER JOIN market_values AS m
	ON mvdq.market_value_id = m.id
	INNER JOIN market_value_diagnose_choices AS mdc
	ON mdc.id = mvda.choice_id
	INNER JOIN market_value_diagnose_question_weights AS mdqw
	ON mdqw.question_id = mvdq.id AND mdqw.choice_id = mdc.id
WHERE mvda.answer_group_id = ?
GROUP BY m.id
ORDER BY m.display_order
`

	// 合計点数を${MarketValueDiagnosePerfectPoint}点満点で換算
	query := fmt.Sprintf(`
SELECT
	tpmv.total_point * %d DIV spmv.sum_point AS total_point,
	tpmv.market_value_name
FROM (%s) AS tpmv
	INNER JOIN (%s) AS spmv
	ON spmv.market_value_id = tpmv.market_value_id
`, MarketValueDiagnosePerfectPoint, totalPointByMarketValueQuery, marketValueSumPointByMarketValueQuery)

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, answerGroupId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customMarketValueDiagnoseAnswersModel) FindAnswerResultList(ctx context.Context, userId string) (results []MarketValueDiagnoseAnswerResult, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAnswerResultList error")
	}()

	var resp []MarketValueDiagnoseAnswerResult

	totalPointByMarketValueQuery := `
SELECT
	SUM(mdqw.weight) AS total_point,
	m.id AS market_value_id,
	m.name AS market_value_name
FROM market_value_diagnose_answers AS mvda
	INNER JOIN market_value_diagnose_questions AS mvdq
	ON mvdq.id = mvda.question_id
	INNER JOIN market_values AS m
	ON mvdq.market_value_id = m.id
	INNER JOIN market_value_diagnose_choices AS mdc
	ON mdc.id = mvda.choice_id
	INNER JOIN market_value_diagnose_question_weights AS mdqw
	ON mdqw.question_id = mvdq.id AND mdqw.choice_id = mdc.id
WHERE mvda.answer_group_id = (
	SELECT
		answer_group_id 
	FROM market_value_diagnose_answers
	WHERE user_id = ?
	ORDER BY answer_group_id DESC
	LIMIT 1
)
GROUP BY m.id
ORDER BY m.display_order
`

	// 合計点数を${MarketValueDiagnosePerfectPoint}点満点で換算
	query := fmt.Sprintf(`
SELECT
	tpmv.total_point * %d DIV spmv.sum_point AS total_point,
	tpmv.market_value_name
FROM (%s) AS tpmv
	INNER JOIN (%s) AS spmv
	ON spmv.market_value_id = tpmv.market_value_id
`, MarketValueDiagnosePerfectPoint, totalPointByMarketValueQuery, marketValueSumPointByMarketValueQuery)

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customMarketValueDiagnoseAnswersModel) FindLatestAnswerList(ctx context.Context, userId string) (answers []MarketValueDiagnoseAnswers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindLatestAnswerList error")
	}()

	query := `
SELECT
	*
FROM market_value_diagnose_answers AS mvda
WHERE mvda.answer_group_id = (
	SELECT
		answer_group_id 
	FROM market_value_diagnose_answers
	WHERE user_id = ?
	ORDER BY answer_group_id DESC
	LIMIT 1
)
`

	err = m.QueryRowsNoCacheCtx(ctx, &answers, query, userId)

	switch err {
	case nil:
		return answers, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
