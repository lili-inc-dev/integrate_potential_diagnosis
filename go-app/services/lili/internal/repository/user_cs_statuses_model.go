package repository

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserCsStatusesModel = (*customUserCsStatusesModel)(nil)

type (
	// UserCsStatusesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCsStatusesModel.
	UserCsStatusesModel interface {
		userCsStatusesModel
		FindOneByUserIdNoCache(ctx context.Context, userId string) (csStatus *UserCsStatuses, err error)
		UpsertNoCache(ctx context.Context, userId string, csStatus string) (err error)
		UpdateStatus(ctx context.Context, id string, csStatus string) (err error)
	}

	customUserCsStatusesModel struct {
		*defaultUserCsStatusesModel
	}
)

const CsStatusRegisterOnly = "register_only"
const CsStatusDiagnosing = "diagnosing"
const CsStatusNormal = "normal"

// NewUserCsStatusesModel returns a model for the database table.
func NewUserCsStatusesModel(conn sqlx.SqlConn, c cache.CacheConf) UserCsStatusesModel {
	return &customUserCsStatusesModel{
		defaultUserCsStatusesModel: newUserCsStatusesModel(conn, c),
	}
}

func (m *customUserCsStatusesModel) FindOneByUserIdNoCache(ctx context.Context, userId string) (csStatus *UserCsStatuses, err error) {
	defer func() {
		err = errors.Wrap(err, "FindOneNoCache error")
	}()

	var resp UserCsStatuses
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userCsStatusesRows, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserCsStatusesModel) UpsertNoCache(ctx context.Context, userId string, csStatus string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpsertNoCache error")
	}()

	ulid, err := util.GenerateUlid()
	if err != nil {
		return err
	}

	query := `
INSERT INTO user_cs_statuses(id, user_id, status)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
status = VALUES(status)
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		ulid.String(),
		userId,
		csStatus,
	)

	return err
}

func (m *customUserCsStatusesModel) UpdateStatus(ctx context.Context, id string, csStatus string) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateStatus error")
	}()

	query := `
UPDATE
user_cs_statuses
SET status = ?
WHERE id = ?
`
	_, err = m.ExecNoCacheCtx(
		ctx,
		query,
		csStatus,
		id,
	)

	return err
}
