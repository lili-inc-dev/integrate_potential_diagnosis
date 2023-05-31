// Code generated by goctl. DO NOT EDIT!

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	marketValueDiagnoseAnswersFieldNames          = builder.RawFieldNames(&MarketValueDiagnoseAnswers{})
	marketValueDiagnoseAnswersRows                = strings.Join(marketValueDiagnoseAnswersFieldNames, ",")
	marketValueDiagnoseAnswersRowsExpectAutoSet   = strings.Join(stringx.Remove(marketValueDiagnoseAnswersFieldNames, "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	marketValueDiagnoseAnswersRowsWithPlaceHolder = strings.Join(stringx.Remove(marketValueDiagnoseAnswersFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheMarketValueDiagnoseAnswersIdPrefix                      = "cache:marketValueDiagnoseAnswers:id:"
	cacheMarketValueDiagnoseAnswersAnswerGroupIdQuestionIdPrefix = "cache:marketValueDiagnoseAnswers:answerGroupId:questionId:"
)

type (
	marketValueDiagnoseAnswersModel interface {
		Insert(ctx context.Context, data *MarketValueDiagnoseAnswers) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*MarketValueDiagnoseAnswers, error)
		FindOneByAnswerGroupIdQuestionId(ctx context.Context, answerGroupId string, questionId uint64) (*MarketValueDiagnoseAnswers, error)
		Update(ctx context.Context, data *MarketValueDiagnoseAnswers) error
		Delete(ctx context.Context, id string) error
	}

	defaultMarketValueDiagnoseAnswersModel struct {
		sqlc.CachedConn
		table string
	}

	MarketValueDiagnoseAnswers struct {
		Id            string    `db:"id"`              // ULID
		AnswerGroupId string    `db:"answer_group_id"` // ULID 回答グループごとに作成
		QuestionId    uint64    `db:"question_id"`     // 設問ID
		ChoiceId      uint64    `db:"choice_id"`       // 選択した選択肢のID
		UserId        string    `db:"user_id"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}
)

func newMarketValueDiagnoseAnswersModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMarketValueDiagnoseAnswersModel {
	return &defaultMarketValueDiagnoseAnswersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`market_value_diagnose_answers`",
	}
}

func (m *defaultMarketValueDiagnoseAnswersModel) Delete(ctx context.Context, id string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseAnswersAnswerGroupIdQuestionIdPrefix, data.AnswerGroupId, data.QuestionId)
	marketValueDiagnoseAnswersIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseAnswersIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey, marketValueDiagnoseAnswersIdKey)
	return err
}

func (m *defaultMarketValueDiagnoseAnswersModel) FindOne(ctx context.Context, id string) (*MarketValueDiagnoseAnswers, error) {
	marketValueDiagnoseAnswersIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseAnswersIdPrefix, id)
	var resp MarketValueDiagnoseAnswers
	err := m.QueryRowCtx(ctx, &resp, marketValueDiagnoseAnswersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseAnswersRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMarketValueDiagnoseAnswersModel) FindOneByAnswerGroupIdQuestionId(ctx context.Context, answerGroupId string, questionId uint64) (*MarketValueDiagnoseAnswers, error) {
	marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseAnswersAnswerGroupIdQuestionIdPrefix, answerGroupId, questionId)
	var resp MarketValueDiagnoseAnswers
	err := m.QueryRowIndexCtx(ctx, &resp, marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `answer_group_id` = ? and `question_id` = ? limit 1", marketValueDiagnoseAnswersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, answerGroupId, questionId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMarketValueDiagnoseAnswersModel) Insert(ctx context.Context, data *MarketValueDiagnoseAnswers) (sql.Result, error) {
	marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseAnswersAnswerGroupIdQuestionIdPrefix, data.AnswerGroupId, data.QuestionId)
	marketValueDiagnoseAnswersIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseAnswersIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, marketValueDiagnoseAnswersRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.AnswerGroupId, data.QuestionId, data.ChoiceId, data.UserId, data.CreatedAt, data.UpdatedAt)
	}, marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey, marketValueDiagnoseAnswersIdKey)
	return ret, err
}

func (m *defaultMarketValueDiagnoseAnswersModel) Update(ctx context.Context, newData *MarketValueDiagnoseAnswers) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseAnswersAnswerGroupIdQuestionIdPrefix, data.AnswerGroupId, data.QuestionId)
	marketValueDiagnoseAnswersIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseAnswersIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, marketValueDiagnoseAnswersRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.AnswerGroupId, newData.QuestionId, newData.ChoiceId, newData.UserId, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	}, marketValueDiagnoseAnswersAnswerGroupIdQuestionIdKey, marketValueDiagnoseAnswersIdKey)
	return err
}

func (m *defaultMarketValueDiagnoseAnswersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheMarketValueDiagnoseAnswersIdPrefix, primary)
}

func (m *defaultMarketValueDiagnoseAnswersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseAnswersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMarketValueDiagnoseAnswersModel) tableName() string {
	return m.table
}