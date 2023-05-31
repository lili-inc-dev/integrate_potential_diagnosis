package repository

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PersonalityDiagnoseAnswersModel = (*customPersonalityDiagnoseAnswersModel)(nil)

type (
	// PersonalityDiagnoseAnswersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPersonalityDiagnoseAnswersModel.
	PersonalityDiagnoseAnswersModel interface {
		personalityDiagnoseAnswersModel
		BulkInsert(answers []PersonalityDiagnoseAnswers) error
		FindAnswerResultListByAnswerGroupId(ctx context.Context, answerGroupId string) ([]PersonalityDiagnoseAnswerResult, error)
		FindAnswerResultList(ctx context.Context, userId string) ([]PersonalityDiagnoseAnswerResult, error)
		FindLatestAnswerList(ctx context.Context, userId string) (answers []PersonalityDiagnoseAnswers, err error)
	}

	customPersonalityDiagnoseAnswersModel struct {
		*defaultPersonalityDiagnoseAnswersModel
		bulkInserter sqlx.BulkInserter
	}

	PersonalityDiagnoseAnswerResult struct {
		PersonalityName string `db:"personality_name"`
		TotalPoint      uint64 `db:"total_point"`
	}
)

const insertSql = `
insert into personality_diagnose_answers 
(id,answer_group_id,question_id,choice_id,user_id)
values (?,?,?,?,?)`

// NewPersonalityDiagnoseAnswersModel returns a model for the database table.
func NewPersonalityDiagnoseAnswersModel(conn sqlx.SqlConn, c cache.CacheConf) PersonalityDiagnoseAnswersModel {
	bulkInserter, _ := sqlx.NewBulkInserter(conn, insertSql)
	return &customPersonalityDiagnoseAnswersModel{
		defaultPersonalityDiagnoseAnswersModel: newPersonalityDiagnoseAnswersModel(conn, c),
		bulkInserter:                           *bulkInserter,
	}
}

func (m *customPersonalityDiagnoseAnswersModel) BulkInsert(answers []PersonalityDiagnoseAnswers) (err error) {
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

func (m *customPersonalityDiagnoseAnswersModel) FindAnswerResultListByAnswerGroupId(ctx context.Context, answerGroupId string) (result []PersonalityDiagnoseAnswerResult, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAnswerResultListByAnswerGroupId error")
	}()

	var resp []PersonalityDiagnoseAnswerResult

	query := `
SELECT
	SUM(pdc.weight) AS total_point,
	p.name AS personality_name
FROM personality_diagnose_answers AS pda
	INNER JOIN personality_diagnose_questions AS pdq
	ON pdq.id = pda.question_id
	INNER JOIN personalities AS p
	ON pdq.personality_id = p.id
	INNER JOIN personality_diagnose_choices AS pdc
	ON pdc.id = pda.choice_id
WHERE pda.answer_group_id = ?
GROUP BY p.id
ORDER BY MIN(pdq.index)
`
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

func (m *customPersonalityDiagnoseAnswersModel) FindAnswerResultList(ctx context.Context, userId string) (result []PersonalityDiagnoseAnswerResult, err error) {
	defer func() {
		err = errors.Wrap(err, "FindAnswerResultList error")
	}()

	var resp []PersonalityDiagnoseAnswerResult

	query := `
SELECT
	SUM(pdc.weight) AS total_point,
	p.name AS personality_name
FROM personality_diagnose_answers AS pda
	INNER JOIN personality_diagnose_questions AS pdq
	ON pdq.id = pda.question_id
	INNER JOIN personalities AS p
	ON pdq.personality_id = p.id
	INNER JOIN personality_diagnose_choices AS pdc
	ON pdc.id = pda.choice_id
WHERE pda.answer_group_id = (
	SELECT
		answer_group_id 
	FROM personality_diagnose_answers
	WHERE user_id = ?
	ORDER BY answer_group_id DESC
	LIMIT 1
)
GROUP BY p.id
ORDER BY MIN(pdq.index)
`
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

func (m *customPersonalityDiagnoseAnswersModel) FindLatestAnswerList(ctx context.Context, userId string) (answers []PersonalityDiagnoseAnswers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindLatestAnswerList error")
	}()

	var resp []PersonalityDiagnoseAnswers

	query := `
SELECT
	*
FROM personality_diagnose_answers AS pda
WHERE pda.answer_group_id = (
	SELECT
		answer_group_id 
	FROM personality_diagnose_answers
	WHERE user_id = ?
	ORDER BY answer_group_id DESC
	LIMIT 1
)
`
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
