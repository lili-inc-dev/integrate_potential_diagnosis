package repository

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketValueDiagnoseQuestionWeightsModel = (*customMarketValueDiagnoseQuestionWeightsModel)(nil)

type (
	// MarketValueDiagnoseQuestionWeightsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketValueDiagnoseQuestionWeightsModel.
	MarketValueDiagnoseQuestionWeightsModel interface {
		marketValueDiagnoseQuestionWeightsModel
	}

	customMarketValueDiagnoseQuestionWeightsModel struct {
		*defaultMarketValueDiagnoseQuestionWeightsModel
	}
)

// NewMarketValueDiagnoseQuestionWeightsModel returns a model for the database table.
func NewMarketValueDiagnoseQuestionWeightsModel(conn sqlx.SqlConn, c cache.CacheConf) MarketValueDiagnoseQuestionWeightsModel {
	return &customMarketValueDiagnoseQuestionWeightsModel{
		defaultMarketValueDiagnoseQuestionWeightsModel: newMarketValueDiagnoseQuestionWeightsModel(conn, c),
	}
}
