// Code generated by goctl. DO NOT EDIT!

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	gendersFieldNames          = builder.RawFieldNames(&Genders{})
	gendersRows                = strings.Join(gendersFieldNames, ",")
	gendersRowsExpectAutoSet   = strings.Join(stringx.Remove(gendersFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	gendersRowsWithPlaceHolder = strings.Join(stringx.Remove(gendersFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheGendersIdPrefix   = "cache:genders:id:"
	cacheGendersNamePrefix = "cache:genders:name:"
)

type (
	gendersModel interface {
		Insert(ctx context.Context, data *Genders) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Genders, error)
		FindOneByName(ctx context.Context, name string) (*Genders, error)
		Update(ctx context.Context, data *Genders) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultGendersModel struct {
		sqlc.CachedConn
		table string
	}

	Genders struct {
		Id   uint64 `db:"id"`
		Name string `db:"name"` // 男性 or 女性
	}
)

func newGendersModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultGendersModel {
	return &defaultGendersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`genders`",
	}
}

func (m *defaultGendersModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	gendersIdKey := fmt.Sprintf("%s%v", cacheGendersIdPrefix, id)
	gendersNameKey := fmt.Sprintf("%s%v", cacheGendersNamePrefix, data.Name)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, gendersIdKey, gendersNameKey)
	return err
}

func (m *defaultGendersModel) FindOne(ctx context.Context, id uint64) (*Genders, error) {
	gendersIdKey := fmt.Sprintf("%s%v", cacheGendersIdPrefix, id)
	var resp Genders
	err := m.QueryRowCtx(ctx, &resp, gendersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gendersRows, m.table)
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

func (m *defaultGendersModel) FindOneByName(ctx context.Context, name string) (*Genders, error) {
	gendersNameKey := fmt.Sprintf("%s%v", cacheGendersNamePrefix, name)
	var resp Genders
	err := m.QueryRowIndexCtx(ctx, &resp, gendersNameKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", gendersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, name); err != nil {
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

func (m *defaultGendersModel) Insert(ctx context.Context, data *Genders) (sql.Result, error) {
	gendersIdKey := fmt.Sprintf("%s%v", cacheGendersIdPrefix, data.Id)
	gendersNameKey := fmt.Sprintf("%s%v", cacheGendersNamePrefix, data.Name)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?)", m.table, gendersRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name)
	}, gendersIdKey, gendersNameKey)
	return ret, err
}

func (m *defaultGendersModel) Update(ctx context.Context, newData *Genders) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	gendersIdKey := fmt.Sprintf("%s%v", cacheGendersIdPrefix, data.Id)
	gendersNameKey := fmt.Sprintf("%s%v", cacheGendersNamePrefix, data.Name)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, gendersRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Name, newData.Id)
	}, gendersIdKey, gendersNameKey)
	return err
}

func (m *defaultGendersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGendersIdPrefix, primary)
}

func (m *defaultGendersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", gendersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultGendersModel) tableName() string {
	return m.table
}
