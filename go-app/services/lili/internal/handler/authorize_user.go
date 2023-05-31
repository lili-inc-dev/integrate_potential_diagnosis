package handler

import (
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/pkg/errors"
	"net/http"
)

func authorizeUser(
	typeModel repository.UserTypesModel,
	r *http.Request,
	check func(userType *repository.UserTypes) bool,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "authorizeUser error")
	}()

	user, ok := r.Context().Value(constant.CtxKeyUser).(*repository.Users)
	if !ok {
		return errors.New("context cast error")
	}

	userType, err := typeModel.FindOne(r.Context(), user.TypeId)
	if err != nil {
		return err
	}

	if !check(userType) {
		err = errors.New("authorization error")
		return errorx.New(err, errorx.Forbidden)
	}

	return nil
}
