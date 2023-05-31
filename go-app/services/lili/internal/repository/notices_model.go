package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ NoticesModel = (*customNoticesModel)(nil)

type (
	// NoticesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticesModel.
	NoticesModel interface {
		noticesModel
		FindOneByIdNoCache(id uint64) (*Notices, error)
		FindSpecificRangeNoticeAll(ctx context.Context, limit uint64, offset uint64) ([]Notices, error)
		FindCountNoticeAll(ctx context.Context) (*uint64, error)
		FindSpecificRangeNotice(ctx context.Context, limit uint64, offset uint64) ([]Notices, error)
		FindCountReleasedNotice(ctx context.Context) (*uint64, error)
		InsertNoCache(ctx context.Context, title, content string, isReleased bool) (sql.Result, error)
		UpdateNoCache(ctx context.Context, id uint64, title, content string, isReleased bool) error
	}

	customNoticesModel struct {
		*defaultNoticesModel
	}
)

// NewNoticesModel returns a model for the database table.
func NewNoticesModel(conn sqlx.SqlConn, c cache.CacheConf) NoticesModel {
	return &customNoticesModel{
		defaultNoticesModel: newNoticesModel(conn, c),
	}
}

func (m *customNoticesModel) FindOneByIdNoCache(id uint64) (notices *Notices, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneByIdNoCache error")
	}()

	var resp Notices
	query := fmt.Sprintf(
		"select %s from %s where id = ?",
		noticesRows,
		m.table,
	)
	err = m.QueryRowNoCache(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customNoticesModel) FindSpecificRangeNoticeAll(ctx context.Context, limit uint64, offset uint64) (notices []Notices, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeNoticeAll error")
	}()

	var resp []Notices
	query := fmt.Sprintf(
		"select %s from %s order by id desc limit %s offset %s",
		noticesRows,
		m.table,
		strconv.FormatUint(limit, 10),
		strconv.FormatUint(offset, 10),
	)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customNoticesModel) FindCountNoticeAll(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCountNoticeAll error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}
	query := fmt.Sprintf(
		"select COUNT(`id`) AS `count` from %s",
		m.table,
	)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return &resp.Count, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customNoticesModel) FindSpecificRangeNotice(ctx context.Context, limit uint64, offset uint64) (notices []Notices, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeNotice error")
	}()

	var resp []Notices
	query := fmt.Sprintf(
		"select %s from %s where `is_released` = true limit %s offset %s",
		noticesRows,
		m.table,
		strconv.FormatUint(limit, 10),
		strconv.FormatUint(offset, 10),
	)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customNoticesModel) FindCountReleasedNotice(ctx context.Context) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCountReleasedNotice error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}
	query := fmt.Sprintf(
		"select COUNT(`id`) AS `count` from %s where `is_released` = true",
		m.table,
	)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query)

	switch err {
	case nil:
		return &resp.Count, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customNoticesModel) InsertNoCache(ctx context.Context, title, content string, isReleased bool) (res sql.Result, err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	rows := strings.Join(stringx.Remove(noticesFieldNames, "`id`", "`created_at`", "`updated_at`"), ",")
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, rows)
	ret, err := m.ExecNoCacheCtx(
		ctx,
		query,
		title,
		content,
		isReleased,
	)

	return ret, err
}

func (m *customNoticesModel) UpdateNoCache(ctx context.Context, id uint64, title, content string, isReleased bool) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateNoCache error")
	}()

	query := `
UPDATE
notices
SET 
	title = ?,
	content = ?,
	is_released = ?
WHERE id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		title,
		content,
		isReleased,
		id,
	)
	return err
}
