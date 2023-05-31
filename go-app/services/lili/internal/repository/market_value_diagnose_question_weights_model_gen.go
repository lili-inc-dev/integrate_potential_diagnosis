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
	marketValueDiagnoseQuestionWeightsFieldNames          = builder.RawFieldNames(&MarketValueDiagnoseQuestionWeights{})
	marketValueDiagnoseQuestionWeightsRows                = strings.Join(marketValueDiagnoseQuestionWeightsFieldNames, ",")
	marketValueDiagnoseQuestionWeightsRowsExpectAutoSet   = strings.Join(stringx.Remove(marketValueDiagnoseQuestionWeightsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	marketValueDiagnoseQuestionWeightsRowsWithPlaceHolder = strings.Join(stringx.Remove(marketValueDiagnoseQuestionWeightsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheMarketValueDiagnoseQuestionWeightsIdPrefix                 = "cache:marketValueDiagnoseQuestionWeights:id:"
	cacheMarketValueDiagnoseQuestionWeightsQuestionIdChoiceIdPrefix = "cache:marketValueDiagnoseQuestionWeights:questionId:choiceId:"
)

type (
	marketValueDiagnoseQuestionWeightsModel interface {
		Insert(ctx context.Context, data *MarketValueDiagnoseQuestionWeights) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*MarketValueDiagnoseQuestionWeights, error)
		FindOneByQuestionIdChoiceId(ctx context.Context, questionId uint64, choiceId uint64) (*MarketValueDiagnoseQuestionWeights, error)
		Update(ctx context.Context, data *MarketValueDiagnoseQuestionWeights) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMarketValueDiagnoseQuestionWeightsModel struct {
		sqlc.CachedConn
		table string
	}

	MarketValueDiagnoseQuestionWeights struct {
		Id         uint64    `db:"id"`
		QuestionId uint64    `db:"question_id"`
		ChoiceId   uint64    `db:"choice_id"`
		Weight     uint64    `db:"weight"` // 回答の重み
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}
)

func newMarketValueDiagnoseQuestionWeightsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMarketValueDiagnoseQuestionWeightsModel {
	return &defaultMarketValueDiagnoseQuestionWeightsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`market_value_diagnose_question_weights`",
	}
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	marketValueDiagnoseQuestionWeightsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionWeightsIdPrefix, id)
	marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseQuestionWeightsQuestionIdChoiceIdPrefix, data.QuestionId, data.ChoiceId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, marketValueDiagnoseQuestionWeightsIdKey, marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey)
	return err
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) FindOne(ctx context.Context, id uint64) (*MarketValueDiagnoseQuestionWeights, error) {
	marketValueDiagnoseQuestionWeightsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionWeightsIdPrefix, id)
	var resp MarketValueDiagnoseQuestionWeights
	err := m.QueryRowCtx(ctx, &resp, marketValueDiagnoseQuestionWeightsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseQuestionWeightsRows, m.table)
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

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) FindOneByQuestionIdChoiceId(ctx context.Context, questionId uint64, choiceId uint64) (*MarketValueDiagnoseQuestionWeights, error) {
	marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseQuestionWeightsQuestionIdChoiceIdPrefix, questionId, choiceId)
	var resp MarketValueDiagnoseQuestionWeights
	err := m.QueryRowIndexCtx(ctx, &resp, marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `question_id` = ? and `choice_id` = ? limit 1", marketValueDiagnoseQuestionWeightsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, questionId, choiceId); err != nil {
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

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) Insert(ctx context.Context, data *MarketValueDiagnoseQuestionWeights) (sql.Result, error) {
	marketValueDiagnoseQuestionWeightsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionWeightsIdPrefix, data.Id)
	marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseQuestionWeightsQuestionIdChoiceIdPrefix, data.QuestionId, data.ChoiceId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, marketValueDiagnoseQuestionWeightsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.QuestionId, data.ChoiceId, data.Weight, data.CreatedAt, data.UpdatedAt)
	}, marketValueDiagnoseQuestionWeightsIdKey, marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey)
	return ret, err
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) Update(ctx context.Context, newData *MarketValueDiagnoseQuestionWeights) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	marketValueDiagnoseQuestionWeightsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionWeightsIdPrefix, data.Id)
	marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey := fmt.Sprintf("%s%v:%v", cacheMarketValueDiagnoseQuestionWeightsQuestionIdChoiceIdPrefix, data.QuestionId, data.ChoiceId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, marketValueDiagnoseQuestionWeightsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.QuestionId, newData.ChoiceId, newData.Weight, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	}, marketValueDiagnoseQuestionWeightsIdKey, marketValueDiagnoseQuestionWeightsQuestionIdChoiceIdKey)
	return err
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionWeightsIdPrefix, primary)
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseQuestionWeightsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMarketValueDiagnoseQuestionWeightsModel) tableName() string {
	return m.table
}
