package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	DiagnosesModel interface {
		FindCountByUserId(ctx context.Context, userId string) (*uint64, error)
		FindSpecificRangeDiagnoseHistoryByUserId(userId string, limit uint64, offset uint64) ([]DiagnoseHistory, error)
	}

	customDiagnosesModel struct {
		sqlc.CachedConn
	}

	DiagnoseHistory struct {
		DiagnoseType  string    `db:"diagnose_type"`
		AnswerGroupId string    `db:"answer_group_id"`
		DiagnoseName  string    `db:"diagnose_name"`
		CreatedAt     time.Time `db:"created_at"`
	}
)

const (
	PersonalityDiagnoseName = "ヒューマンリテラシー"
	MarketValueDiagnoseName = "ビジネスリテラシー"
	CareerWorkName          = "キャリアワーク"
)

// NewDiagnosesModel returns a model for the database table.
func NewDiagnosesModel(conn sqlx.SqlConn, c cache.CacheConf) DiagnosesModel {
	return &customDiagnosesModel{
		CachedConn: sqlc.NewConn(conn, c),
	}
}

func (m *customDiagnosesModel) FindCountByUserId(ctx context.Context, userId string) (count *uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "FindCountByUserId error")
	}()

	resp := struct {
		Count uint64 `db:"count"`
	}{}

	query := `
SELECT
 COUNT(*) AS count
FROM ((
SELECT
	pda.answer_group_id
FROM personality_diagnose_answers AS pda
WHERE pda.user_id = ?
GROUP BY answer_group_id
)
UNION ALL
(
SELECT
	mvda.answer_group_id
FROM market_value_diagnose_answers AS mvda
WHERE mvda.user_id = ?
GROUP BY answer_group_id
)
UNION ALL
(
SELECT
	cwa.answer_group_id
FROM career_work_answers AS cwa
WHERE cwa.user_id = ?
GROUP BY answer_group_id
)) as diagnose
`

	err = m.QueryRowNoCacheCtx(ctx, &resp, query, userId, userId, userId)

	switch err {
	case nil:
		return &resp.Count, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customDiagnosesModel) FindSpecificRangeDiagnoseHistoryByUserId(userId string, limit uint64, offset uint64) (histories []DiagnoseHistory, err error) {
	defer func() {
		err = errors.Wrap(err, "FindSpecificRangeDiagnoseHistoryByUserId error")
	}()

	var resp []DiagnoseHistory
	query := fmt.Sprintf(`
(
SELECT
	pda.answer_group_id,
	'personality' AS diagnose_type,
	'%s' AS diagnose_name,
	MIN(pda.created_at) AS created_at 
FROM personality_diagnose_answers AS pda
WHERE pda.user_id = ?
GROUP BY answer_group_id
)
UNION ALL
(
SELECT
	mvda.answer_group_id,
	'market_value' AS diagnose_type,
	'%s' AS diagnose_name,
	MIN(mvda.created_at) AS created_at 
FROM market_value_diagnose_answers AS mvda
WHERE mvda.user_id = ?
GROUP BY answer_group_id
)
UNION ALL
(
SELECT
	cwa.answer_group_id,
	'career_work' AS diagnose_type,
	'%s' AS diagnose_name,
	MIN(cwa.created_at) AS created_at 
FROM career_work_answers AS cwa
WHERE cwa.user_id = ?
GROUP BY answer_group_id
)
ORDER BY answer_group_id DESC
LIMIT ?
OFFSET ?
`, PersonalityDiagnoseName, MarketValueDiagnoseName, CareerWorkName)

	err = m.QueryRowsNoCache(&resp, query, userId, userId, userId, limit, offset)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
