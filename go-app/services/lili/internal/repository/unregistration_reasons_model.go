package repository

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UnregistrationReasonsModel = (*customUnregistrationReasonsModel)(nil)

type (
	// UnregistrationReasonsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUnregistrationReasonsModel.
	UnregistrationReasonsModel interface {
		unregistrationReasonsModel
	}

	customUnregistrationReasonsModel struct {
		*defaultUnregistrationReasonsModel
	}
)

// NewUnregistrationReasonsModel returns a model for the database table.
func NewUnregistrationReasonsModel(conn sqlx.SqlConn, c cache.CacheConf) UnregistrationReasonsModel {
	return &customUnregistrationReasonsModel{
		defaultUnregistrationReasonsModel: newUnregistrationReasonsModel(conn, c),
	}
}
