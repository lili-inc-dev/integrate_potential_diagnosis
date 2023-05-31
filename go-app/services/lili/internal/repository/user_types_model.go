package repository

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserTypesModel = (*customUserTypesModel)(nil)

type (
	// UserTypesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTypesModel.
	UserTypesModel interface {
		userTypesModel
	}

	customUserTypesModel struct {
		*defaultUserTypesModel
	}
)

const (
	UserTypeNameGeneral = "学生"
	UserTypeNameCompany = "社会人"
)

// NewUserTypesModel returns a model for the database table.
func NewUserTypesModel(conn sqlx.SqlConn, c cache.CacheConf) UserTypesModel {
	return &customUserTypesModel{
		defaultUserTypesModel: newUserTypesModel(conn, c),
	}
}
