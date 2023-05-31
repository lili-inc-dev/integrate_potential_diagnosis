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
	unregisteredUserProfilesFieldNames          = builder.RawFieldNames(&UnregisteredUserProfiles{})
	unregisteredUserProfilesRows                = strings.Join(unregisteredUserProfilesFieldNames, ",")
	unregisteredUserProfilesRowsExpectAutoSet   = strings.Join(stringx.Remove(unregisteredUserProfilesFieldNames, "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	unregisteredUserProfilesRowsWithPlaceHolder = strings.Join(stringx.Remove(unregisteredUserProfilesFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUnregisteredUserProfilesIdPrefix     = "cache:unregisteredUserProfiles:id:"
	cacheUnregisteredUserProfilesUserIdPrefix = "cache:unregisteredUserProfiles:userId:"
)

type (
	unregisteredUserProfilesModel interface {
		Insert(ctx context.Context, data *UnregisteredUserProfiles) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*UnregisteredUserProfiles, error)
		FindOneByUserId(ctx context.Context, userId string) (*UnregisteredUserProfiles, error)
		Update(ctx context.Context, data *UnregisteredUserProfiles) error
		Delete(ctx context.Context, id string) error
	}

	defaultUnregisteredUserProfilesModel struct {
		sqlc.CachedConn
		table string
	}

	UnregisteredUserProfiles struct {
		Id        string    `db:"id"` // ULID
		UserId    string    `db:"user_id"`
		ReasonId  uint64    `db:"reason_id"`
		LineId    string    `db:"line_id"`
		Email     string    `db:"email"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func newUnregisteredUserProfilesModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUnregisteredUserProfilesModel {
	return &defaultUnregisteredUserProfilesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`unregistered_user_profiles`",
	}
}

func (m *defaultUnregisteredUserProfilesModel) Delete(ctx context.Context, id string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	unregisteredUserProfilesIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesIdPrefix, id)
	unregisteredUserProfilesUserIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesUserIdPrefix, data.UserId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, unregisteredUserProfilesIdKey, unregisteredUserProfilesUserIdKey)
	return err
}

func (m *defaultUnregisteredUserProfilesModel) FindOne(ctx context.Context, id string) (*UnregisteredUserProfiles, error) {
	unregisteredUserProfilesIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesIdPrefix, id)
	var resp UnregisteredUserProfiles
	err := m.QueryRowCtx(ctx, &resp, unregisteredUserProfilesIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", unregisteredUserProfilesRows, m.table)
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

func (m *defaultUnregisteredUserProfilesModel) FindOneByUserId(ctx context.Context, userId string) (*UnregisteredUserProfiles, error) {
	unregisteredUserProfilesUserIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesUserIdPrefix, userId)
	var resp UnregisteredUserProfiles
	err := m.QueryRowIndexCtx(ctx, &resp, unregisteredUserProfilesUserIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", unregisteredUserProfilesRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId); err != nil {
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

func (m *defaultUnregisteredUserProfilesModel) Insert(ctx context.Context, data *UnregisteredUserProfiles) (sql.Result, error) {
	unregisteredUserProfilesIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesIdPrefix, data.Id)
	unregisteredUserProfilesUserIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, unregisteredUserProfilesRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.ReasonId, data.LineId, data.Email, data.CreatedAt, data.UpdatedAt)
	}, unregisteredUserProfilesIdKey, unregisteredUserProfilesUserIdKey)
	return ret, err
}

func (m *defaultUnregisteredUserProfilesModel) Update(ctx context.Context, newData *UnregisteredUserProfiles) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	unregisteredUserProfilesIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesIdPrefix, data.Id)
	unregisteredUserProfilesUserIdKey := fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesUserIdPrefix, data.UserId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, unregisteredUserProfilesRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.UserId, newData.ReasonId, newData.LineId, newData.Email, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	}, unregisteredUserProfilesIdKey, unregisteredUserProfilesUserIdKey)
	return err
}

func (m *defaultUnregisteredUserProfilesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUnregisteredUserProfilesIdPrefix, primary)
}

func (m *defaultUnregisteredUserProfilesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", unregisteredUserProfilesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUnregisteredUserProfilesModel) tableName() string {
	return m.table
}