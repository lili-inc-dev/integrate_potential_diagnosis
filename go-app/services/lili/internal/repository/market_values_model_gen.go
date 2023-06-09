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
	marketValuesFieldNames          = builder.RawFieldNames(&MarketValues{})
	marketValuesRows                = strings.Join(marketValuesFieldNames, ",")
	marketValuesRowsExpectAutoSet   = strings.Join(stringx.Remove(marketValuesFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	marketValuesRowsWithPlaceHolder = strings.Join(stringx.Remove(marketValuesFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheMarketValuesIdPrefix           = "cache:marketValues:id:"
	cacheMarketValuesDisplayOrderPrefix = "cache:marketValues:displayOrder:"
	cacheMarketValuesNamePrefix         = "cache:marketValues:name:"
)

type (
	marketValuesModel interface {
		Insert(ctx context.Context, data *MarketValues) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*MarketValues, error)
		FindOneByDisplayOrder(ctx context.Context, displayOrder uint64) (*MarketValues, error)
		FindOneByName(ctx context.Context, name string) (*MarketValues, error)
		Update(ctx context.Context, data *MarketValues) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMarketValuesModel struct {
		sqlc.CachedConn
		table string
	}

	MarketValues struct {
		Id           uint64    `db:"id"`
		DisplayOrder uint64    `db:"display_order"` // 何ページ目に表示されるか
		Name         string    `db:"name"`          // 市場価値診断の項目名
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
	}
)

func newMarketValuesModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultMarketValuesModel {
	return &defaultMarketValuesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`market_values`",
	}
}

func (m *defaultMarketValuesModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	marketValuesDisplayOrderKey := fmt.Sprintf("%s%v", cacheMarketValuesDisplayOrderPrefix, data.DisplayOrder)
	marketValuesIdKey := fmt.Sprintf("%s%v", cacheMarketValuesIdPrefix, id)
	marketValuesNameKey := fmt.Sprintf("%s%v", cacheMarketValuesNamePrefix, data.Name)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, marketValuesDisplayOrderKey, marketValuesIdKey, marketValuesNameKey)
	return err
}

func (m *defaultMarketValuesModel) FindOne(ctx context.Context, id uint64) (*MarketValues, error) {
	marketValuesIdKey := fmt.Sprintf("%s%v", cacheMarketValuesIdPrefix, id)
	var resp MarketValues
	err := m.QueryRowCtx(ctx, &resp, marketValuesIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValuesRows, m.table)
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

func (m *defaultMarketValuesModel) FindOneByDisplayOrder(ctx context.Context, displayOrder uint64) (*MarketValues, error) {
	marketValuesDisplayOrderKey := fmt.Sprintf("%s%v", cacheMarketValuesDisplayOrderPrefix, displayOrder)
	var resp MarketValues
	err := m.QueryRowIndexCtx(ctx, &resp, marketValuesDisplayOrderKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `display_order` = ? limit 1", marketValuesRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, displayOrder); err != nil {
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

func (m *defaultMarketValuesModel) FindOneByName(ctx context.Context, name string) (*MarketValues, error) {
	marketValuesNameKey := fmt.Sprintf("%s%v", cacheMarketValuesNamePrefix, name)
	var resp MarketValues
	err := m.QueryRowIndexCtx(ctx, &resp, marketValuesNameKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", marketValuesRows, m.table)
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

func (m *defaultMarketValuesModel) Insert(ctx context.Context, data *MarketValues) (sql.Result, error) {
	marketValuesDisplayOrderKey := fmt.Sprintf("%s%v", cacheMarketValuesDisplayOrderPrefix, data.DisplayOrder)
	marketValuesIdKey := fmt.Sprintf("%s%v", cacheMarketValuesIdPrefix, data.Id)
	marketValuesNameKey := fmt.Sprintf("%s%v", cacheMarketValuesNamePrefix, data.Name)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, marketValuesRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.DisplayOrder, data.Name, data.CreatedAt, data.UpdatedAt)
	}, marketValuesDisplayOrderKey, marketValuesIdKey, marketValuesNameKey)
	return ret, err
}

func (m *defaultMarketValuesModel) Update(ctx context.Context, newData *MarketValues) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	marketValuesDisplayOrderKey := fmt.Sprintf("%s%v", cacheMarketValuesDisplayOrderPrefix, data.DisplayOrder)
	marketValuesIdKey := fmt.Sprintf("%s%v", cacheMarketValuesIdPrefix, data.Id)
	marketValuesNameKey := fmt.Sprintf("%s%v", cacheMarketValuesNamePrefix, data.Name)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, marketValuesRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.DisplayOrder, newData.Name, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	}, marketValuesDisplayOrderKey, marketValuesIdKey, marketValuesNameKey)
	return err
}

func (m *defaultMarketValuesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheMarketValuesIdPrefix, primary)
}

func (m *defaultMarketValuesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", marketValuesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMarketValuesModel) tableName() string {
	return m.table
}
