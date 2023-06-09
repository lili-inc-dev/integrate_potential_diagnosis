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
	marketValueDiagnoseQuestionsFieldNames          = builder.RawFieldNames(&MarketValueDiagnoseQuestions{})
	marketValueDiagnoseQuestionsRows                = strings.Join(marketValueDiagnoseQuestionsFieldNames, ",")
	marketValueDiagnoseQuestionsRowsExpectAutoSet   = strings.Join(stringx.Remove(marketValueDiagnoseQuestionsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	marketValueDiagnoseQuestionsRowsWithPlaceHolder = strings.Join(stringx.Remove(marketValueDiagnoseQuestionsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheMarketValueDiagnoseQuestionsIdPrefix      = "cache:marketValueDiagnoseQuestions:id:"
	cacheMarketValueDiagnoseQuestionsContentPrefix = "cache:marketValueDiagnoseQuestions:content:"
	cacheMarketValueDiagnoseQuestionsIndexPrefix   = "cache:marketValueDiagnoseQuestions:index:"
)

type (
	marketValueDiagnoseQuestionsModel interface {
		Insert(ctx context.Context, data *MarketValueDiagnoseQuestions) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*MarketValueDiagnoseQuestions, error)
		FindOneByContent(ctx context.Context, content string) (*MarketValueDiagnoseQuestions, error)
		FindOneByIndex(ctx context.Context, index int64) (*MarketValueDiagnoseQuestions, error)
		Update(ctx context.Context, data *MarketValueDiagnoseQuestions) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMarketValueDiagnoseQuestionsModel struct {
		sqlc.CachedConn
		table string
	}

	MarketValueDiagnoseQuestions struct {
		Id            uint64    `db:"id"`
		MarketValueId uint64    `db:"market_value_id"` // 市場価値の項目ID
		Index         int64     `db:"index"`           // 出題順
		Content       string    `db:"content"`         // 設問内容
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}
)

func newMarketValueDiagnoseQuestionsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMarketValueDiagnoseQuestionsModel {
	return &defaultMarketValueDiagnoseQuestionsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`market_value_diagnose_questions`",
	}
}

func (m *defaultMarketValueDiagnoseQuestionsModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	marketValueDiagnoseQuestionsContentKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsContentPrefix, data.Content)
	marketValueDiagnoseQuestionsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIdPrefix, id)
	marketValueDiagnoseQuestionsIndexKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIndexPrefix, data.Index)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, marketValueDiagnoseQuestionsContentKey, marketValueDiagnoseQuestionsIdKey, marketValueDiagnoseQuestionsIndexKey)
	return err
}

func (m *defaultMarketValueDiagnoseQuestionsModel) FindOne(ctx context.Context, id uint64) (*MarketValueDiagnoseQuestions, error) {
	marketValueDiagnoseQuestionsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIdPrefix, id)
	var resp MarketValueDiagnoseQuestions
	err := m.QueryRowCtx(ctx, &resp, marketValueDiagnoseQuestionsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseQuestionsRows, m.table)
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

func (m *defaultMarketValueDiagnoseQuestionsModel) FindOneByContent(ctx context.Context, content string) (*MarketValueDiagnoseQuestions, error) {
	marketValueDiagnoseQuestionsContentKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsContentPrefix, content)
	var resp MarketValueDiagnoseQuestions
	err := m.QueryRowIndexCtx(ctx, &resp, marketValueDiagnoseQuestionsContentKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `content` = ? limit 1", marketValueDiagnoseQuestionsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, content); err != nil {
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

func (m *defaultMarketValueDiagnoseQuestionsModel) FindOneByIndex(ctx context.Context, index int64) (*MarketValueDiagnoseQuestions, error) {
	marketValueDiagnoseQuestionsIndexKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIndexPrefix, index)
	var resp MarketValueDiagnoseQuestions
	err := m.QueryRowIndexCtx(ctx, &resp, marketValueDiagnoseQuestionsIndexKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `index` = ? limit 1", marketValueDiagnoseQuestionsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, index); err != nil {
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

func (m *defaultMarketValueDiagnoseQuestionsModel) Insert(ctx context.Context, data *MarketValueDiagnoseQuestions) (sql.Result, error) {
	marketValueDiagnoseQuestionsContentKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsContentPrefix, data.Content)
	marketValueDiagnoseQuestionsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIdPrefix, data.Id)
	marketValueDiagnoseQuestionsIndexKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIndexPrefix, data.Index)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, marketValueDiagnoseQuestionsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.MarketValueId, data.Index, data.Content, data.CreatedAt, data.UpdatedAt)
	}, marketValueDiagnoseQuestionsContentKey, marketValueDiagnoseQuestionsIdKey, marketValueDiagnoseQuestionsIndexKey)
	return ret, err
}

func (m *defaultMarketValueDiagnoseQuestionsModel) Update(ctx context.Context, newData *MarketValueDiagnoseQuestions) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	marketValueDiagnoseQuestionsContentKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsContentPrefix, data.Content)
	marketValueDiagnoseQuestionsIdKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIdPrefix, data.Id)
	marketValueDiagnoseQuestionsIndexKey := fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIndexPrefix, data.Index)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, marketValueDiagnoseQuestionsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.MarketValueId, newData.Index, newData.Content, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	}, marketValueDiagnoseQuestionsContentKey, marketValueDiagnoseQuestionsIdKey, marketValueDiagnoseQuestionsIndexKey)
	return err
}

func (m *defaultMarketValueDiagnoseQuestionsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheMarketValueDiagnoseQuestionsIdPrefix, primary)
}

func (m *defaultMarketValueDiagnoseQuestionsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValueDiagnoseQuestionsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMarketValueDiagnoseQuestionsModel) tableName() string {
	return m.table
}
