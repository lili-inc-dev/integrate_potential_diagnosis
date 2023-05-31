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
	unregistrationReasonsFieldNames          = builder.RawFieldNames(&UnregistrationReasons{})
	unregistrationReasonsRows                = strings.Join(unregistrationReasonsFieldNames, ",")
	unregistrationReasonsRowsExpectAutoSet   = strings.Join(stringx.Remove(unregistrationReasonsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	unregistrationReasonsRowsWithPlaceHolder = strings.Join(stringx.Remove(unregistrationReasonsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUnregistrationReasonsIdPrefix = "cache:unregistrationReasons:id:"
)

type (
	unregistrationReasonsModel interface {
		Insert(ctx context.Context, data *UnregistrationReasons) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*UnregistrationReasons, error)
		Update(ctx context.Context, data *UnregistrationReasons) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultUnregistrationReasonsModel struct {
		sqlc.CachedConn
		table string
	}

	UnregistrationReasons struct {
		Id        uint64    `db:"id"`
		Content   string    `db:"content"` // 退会理由内容
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func newUnregistrationReasonsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUnregistrationReasonsModel {
	return &defaultUnregistrationReasonsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`unregistration_reasons`",
	}
}

func (m *defaultUnregistrationReasonsModel) Delete(ctx context.Context, id uint64) error {
	unregistrationReasonsIdKey := fmt.Sprintf("%s%v", cacheUnregistrationReasonsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, unregistrationReasonsIdKey)
	return err
}

func (m *defaultUnregistrationReasonsModel) FindOne(ctx context.Context, id uint64) (*UnregistrationReasons, error) {
	unregistrationReasonsIdKey := fmt.Sprintf("%s%v", cacheUnregistrationReasonsIdPrefix, id)
	var resp UnregistrationReasons
	err := m.QueryRowCtx(ctx, &resp, unregistrationReasonsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", unregistrationReasonsRows, m.table)
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

func (m *defaultUnregistrationReasonsModel) Insert(ctx context.Context, data *UnregistrationReasons) (sql.Result, error) {
	unregistrationReasonsIdKey := fmt.Sprintf("%s%v", cacheUnregistrationReasonsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, unregistrationReasonsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Content, data.CreatedAt, data.UpdatedAt)
	}, unregistrationReasonsIdKey)
	return ret, err
}

func (m *defaultUnregistrationReasonsModel) Update(ctx context.Context, data *UnregistrationReasons) error {
	unregistrationReasonsIdKey := fmt.Sprintf("%s%v", cacheUnregistrationReasonsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, unregistrationReasonsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Content, data.CreatedAt, data.UpdatedAt, data.Id)
	}, unregistrationReasonsIdKey)
	return err
}

func (m *defaultUnregistrationReasonsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUnregistrationReasonsIdPrefix, primary)
}

func (m *defaultUnregistrationReasonsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", unregistrationReasonsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUnregistrationReasonsModel) tableName() string {
	return m.table
}