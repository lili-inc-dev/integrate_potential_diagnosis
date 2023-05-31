package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAddressesModel = (*customUserAddressesModel)(nil)

type (
	// UserAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAddressesModel.
	UserAddressesModel interface {
		userAddressesModel
		FindByUserIdNoCache(userId string) (*UserAddresses, error)
		InsertNoCache(ctx context.Context, id, userId, postCode, address string) error
		UpdateNoCache(ctx context.Context, userId string, postCode string, address string) error
	}

	customUserAddressesModel struct {
		*defaultUserAddressesModel
	}
)

// NewUserAddressesModel returns a model for the database table.
func NewUserAddressesModel(conn sqlx.SqlConn, c cache.CacheConf) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: newUserAddressesModel(conn, c),
	}
}

func (m *customUserAddressesModel) FindByUserIdNoCache(userId string) (addrese *UserAddresses, err error) {
	defer func() {
		err = errors.Wrap(err, "FindByUserIdNoCache error")
	}()

	var resp UserAddresses
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ?", userAddressesRows, m.table)
	err = m.QueryRowNoCache(&resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserAddressesModel) InsertNoCache(ctx context.Context, id, userId, postCode, address string) (err error) {
	defer func() {
		err = errors.Wrap(err, "InsertNoCache error")
	}()

	query := `
INSERT INTO user_addresses(id, user_id, postal_code, address)
VALUES (?, ?, ?, ?)
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		id,
		userId,
		postCode,
		address,
	)
	return err
}

func (m *customUserAddressesModel) UpdateNoCache(ctx context.Context, userId string, postCode string, address string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateNoCache error")
	}()

	query := `
UPDATE
user_addresses
SET postal_code = ?, address = ?
WHERE user_id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		postCode,
		address,
		userId,
	)
	return err
}
