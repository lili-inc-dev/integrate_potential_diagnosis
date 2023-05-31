package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserInterestTopicsModel = (*customUserInterestTopicsModel)(nil)

type (
	// UserInterestTopicsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInterestTopicsModel.
	UserInterestTopicsModel interface {
		userInterestTopicsModel
		FindListByUserId(ctx context.Context, userId string) ([]UserInterestTopicWithName, error)
		BulkInsert(userId string, topicIds []uint64) error
		DeleteByUserId(userId string) error
	}

	customUserInterestTopicsModel struct {
		*defaultUserInterestTopicsModel
		bulkInserter sqlx.BulkInserter
	}

	UserInterestTopicWithName struct {
		TopicId uint64 `db:"topic_id"`
		Name    string `db:"name"`
		UserId  string `db:"user_id"`
	}
)

const userInterestInsertSql = `
insert ignore into user_interest_topics 
(topic_id,user_id)
values (?,?)`

// NewUserInterestTopicsModel returns a model for the database table.
func NewUserInterestTopicsModel(conn sqlx.SqlConn, c cache.CacheConf) UserInterestTopicsModel {
	bulkInserter, _ := sqlx.NewBulkInserter(conn, userInterestInsertSql)
	return &customUserInterestTopicsModel{
		defaultUserInterestTopicsModel: newUserInterestTopicsModel(conn, c),
		bulkInserter:                   *bulkInserter,
	}
}

func (m *customUserInterestTopicsModel) FindListByUserId(ctx context.Context, userId string) (topic []UserInterestTopicWithName, err error) {
	defer func() {
		err = errors.Wrap(err, "FindListByUserId error")
	}()

	var resp []UserInterestTopicWithName
	query := `
SELECT
	uit.topic_id,
	uit.user_id,
	it.name
FROM user_interest_topics AS uit
	INNER JOIN interest_topics AS it
	ON it.id = uit.topic_id
WHERE uit.user_id = ?
ORDER BY it.display_order
`
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserInterestTopicsModel) BulkInsert(userId string, topicIds []uint64) (err error) {
	defer func() {
		err = errors.Wrap(err, "BulkInsert error")
	}()

	for _, topicId := range topicIds {
		err := m.bulkInserter.Insert(topicId, userId)
		if err != nil {
			return err
		}
	}
	m.bulkInserter.Flush()
	return nil
}

func (m *customUserInterestTopicsModel) DeleteByUserId(userId string) (err error) {
	defer func() {
		err = errors.Wrap(err, "DeleteByUserId error")
	}()

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", m.table)

	_, err = m.ExecNoCache(query, userId)
	return err
}
