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
	usersFieldNames          = builder.RawFieldNames(&Users{})
	usersRows                = strings.Join(usersFieldNames, ",")
	usersRowsExpectAutoSet   = strings.Join(stringx.Remove(usersFieldNames, "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	usersRowsWithPlaceHolder = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUsersIdPrefix          = "cache:users:id:"
	cacheUsersFirebaseUidPrefix = "cache:users:firebaseUid:"
	cacheUsersLineIdPrefix      = "cache:users:lineId:"
)

type (
	usersModel interface {
		Insert(ctx context.Context, data *Users) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*Users, error)
		FindOneByFirebaseUid(ctx context.Context, firebaseUid string) (*Users, error)
		FindOneByLineId(ctx context.Context, lineId sql.NullString) (*Users, error)
		Update(ctx context.Context, data *Users) error
		Delete(ctx context.Context, id string) error
	}

	defaultUsersModel struct {
		sqlc.CachedConn
		table string
	}

	Users struct {
		Id           string         `db:"id"`             // ULID
		LineId       sql.NullString `db:"line_id"`        // LINEユーザID
		TypeId       uint64         `db:"type_id"`        // ユーザ種別
		FirebaseUid  string         `db:"firebase_uid"`   // firebase user id
		Name         string         `db:"name"`           // 氏名
		NameKana     sql.NullString `db:"name_kana"`      // 氏名フリガナ
		GenderId     uint64         `db:"gender_id"`      // 性別
		Birthday     time.Time      `db:"birthday"`       // 生年月日
		PhoneNumber  string         `db:"phone_number"`   // 電話番号
		Memo         sql.NullString `db:"memo"`           // メモ（管理画面で利用）
		LastAccessAt time.Time      `db:"last_access_at"` // 最終アクセス日
		CreatedAt    time.Time      `db:"created_at"`
		UpdatedAt    time.Time      `db:"updated_at"`
		Status       string         `db:"status"` // ステータス（例：registered、banned、unregistered）
	}
)

func newUsersModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUsersModel {
	return &defaultUsersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`users`",
	}
}

func (m *defaultUsersModel) Delete(ctx context.Context, id string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	usersFirebaseUidKey := fmt.Sprintf("%s%v", cacheUsersFirebaseUidPrefix, data.FirebaseUid)
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
	usersLineIdKey := fmt.Sprintf("%s%v", cacheUsersLineIdPrefix, data.LineId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, usersFirebaseUidKey, usersIdKey, usersLineIdKey)
	return err
}

func (m *defaultUsersModel) FindOne(ctx context.Context, id string) (*Users, error) {
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, id)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, usersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
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

func (m *defaultUsersModel) FindOneByFirebaseUid(ctx context.Context, firebaseUid string) (*Users, error) {
	usersFirebaseUidKey := fmt.Sprintf("%s%v", cacheUsersFirebaseUidPrefix, firebaseUid)
	var resp Users
	err := m.QueryRowIndexCtx(ctx, &resp, usersFirebaseUidKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `firebase_uid` = ? limit 1", usersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, firebaseUid); err != nil {
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

func (m *defaultUsersModel) FindOneByLineId(ctx context.Context, lineId sql.NullString) (*Users, error) {
	usersLineIdKey := fmt.Sprintf("%s%v", cacheUsersLineIdPrefix, lineId)
	var resp Users
	err := m.QueryRowIndexCtx(ctx, &resp, usersLineIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `line_id` = ? limit 1", usersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, lineId); err != nil {
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

func (m *defaultUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	usersFirebaseUidKey := fmt.Sprintf("%s%v", cacheUsersFirebaseUidPrefix, data.FirebaseUid)
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, data.Id)
	usersLineIdKey := fmt.Sprintf("%s%v", cacheUsersLineIdPrefix, data.LineId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.LineId, data.TypeId, data.FirebaseUid, data.Name, data.NameKana, data.GenderId, data.Birthday, data.PhoneNumber, data.Memo, data.LastAccessAt, data.CreatedAt, data.UpdatedAt, data.Status)
	}, usersFirebaseUidKey, usersIdKey, usersLineIdKey)
	return ret, err
}

func (m *defaultUsersModel) Update(ctx context.Context, newData *Users) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	usersFirebaseUidKey := fmt.Sprintf("%s%v", cacheUsersFirebaseUidPrefix, data.FirebaseUid)
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, data.Id)
	usersLineIdKey := fmt.Sprintf("%s%v", cacheUsersLineIdPrefix, data.LineId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.LineId, newData.TypeId, newData.FirebaseUid, newData.Name, newData.NameKana, newData.GenderId, newData.Birthday, newData.PhoneNumber, newData.Memo, newData.LastAccessAt, newData.CreatedAt, newData.UpdatedAt, newData.Status, newData.Id)
	}, usersFirebaseUidKey, usersIdKey, usersLineIdKey)
	return err
}

func (m *defaultUsersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUsersIdPrefix, primary)
}

func (m *defaultUsersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUsersModel) tableName() string {
	return m.table
}
