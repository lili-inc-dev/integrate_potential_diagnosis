package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CareerWorkAnswersModel = (*customCareerWorkAnswersModel)(nil)

type (
	// CareerWorkAnswersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCareerWorkAnswersModel.
	CareerWorkAnswersModel interface {
		careerWorkAnswersModel
		FetchListByUserId(ctx context.Context, userId string) ([]CareerWorkAnswers, error)
		FindListByAnswerGroupId(ctx context.Context, answerGroupId string) ([]CareerWorkAnswers, error)
		BulkInsert(answers []CareerWorkAnswers) error
	}

	customCareerWorkAnswersModel struct {
		*defaultCareerWorkAnswersModel
		bulkInserter sqlx.BulkInserter
	}
)

// NewCareerWorkAnswersModel returns a model for the database table.
func NewCareerWorkAnswersModel(conn sqlx.SqlConn, c cache.CacheConf) CareerWorkAnswersModel {
	insertSql := "insert into career_work_answers (`id`, `answer_group_id`, `question_key`, `answer`, `user_id`, `index`) values (?,?,?,?,?,?)"
	bulkInserter, _ := sqlx.NewBulkInserter(conn, insertSql)
	return &customCareerWorkAnswersModel{
		defaultCareerWorkAnswersModel: newCareerWorkAnswersModel(conn, c),
		bulkInserter:                  *bulkInserter,
	}
}

func (m *customCareerWorkAnswersModel) FindListByAnswerGroupId(ctx context.Context, answerGroupId string) (answers []CareerWorkAnswers, err error) {
	defer func() {
		err = errors.Wrap(err, "FindListByAnswerGroupId error")
	}()

	query := fmt.Sprintf(`
SELECT
%s
FROM %s
WHERE answer_group_id = ?
`, careerWorkAnswersRows, m.table)

	err = m.QueryRowsNoCacheCtx(ctx, &answers, query, answerGroupId)

	switch err {
	case nil:
		return answers, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customCareerWorkAnswersModel) FetchListByUserId(ctx context.Context, userId string) (answers []CareerWorkAnswers, err error) {
	defer func() {
		err = errors.Wrap(err, "FetchListByUserId error")
	}()

	query := `
SELECT
	cwa.id,
	cwa.answer_group_id,
	cwa.question_key,
	cwa.answer,
	cwa.user_id,
	cwa.index,
	cwa.created_at,
	cwa.updated_at
FROM career_work_answers AS cwa
WHERE cwa.answer_group_id = (
	SELECT
		answer_group_id 
	FROM career_work_answers
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

func (m *customCareerWorkAnswersModel) BulkInsert(answers []CareerWorkAnswers) (err error) {
	defer func() {
		err = errors.Wrap(err, "BulkInsert error")
	}()

	answerGroupId, err := util.GenerateUlid()
	if err != nil {
		return err
	}
	for _, answer := range answers {
		id, err := util.GenerateUlid()
		if err != nil {
			return err
		}
		err = m.bulkInserter.Insert(id.String(), answerGroupId.String(), answer.QuestionKey, answer.Answer, answer.UserId, answer.Index)
		if err != nil {
			return err
		}
	}
	m.bulkInserter.Flush()
	return nil
}
