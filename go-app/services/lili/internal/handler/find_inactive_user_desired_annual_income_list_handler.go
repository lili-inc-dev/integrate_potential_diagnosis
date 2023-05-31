package handler

import (
	"net/http"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/logic"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/middleware"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindInactiveUserDesiredAnnualIncomeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := middleware.GetFirebaseTokenFromReq(svcCtx.FirebaseAuth, r); err != nil {
			err = errorx.New(err, errorx.Unauthorized)
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserDesiredAnnualIncomeListHandler error"))
			return
		}

		l := logic.NewFindInactiveUserDesiredAnnualIncomeListLogic(r.Context(), svcCtx)
		resp, err := l.FindInactiveUserDesiredAnnualIncomeList()
		if err != nil {
			httpx.Error(w, errors.Wrap(err, "FindInactiveUserDesiredAnnualIncomeListHandler error"))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
