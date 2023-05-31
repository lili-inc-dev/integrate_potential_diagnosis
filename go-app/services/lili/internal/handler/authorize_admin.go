package handler

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
)

func authorizeAdmin(
	roleModel repository.AdminRolesModel,
	r *http.Request,
	check func(role *repository.AdminRoles) bool,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "authorizeAdmin error")
	}()

	admin, ok := r.Context().Value(constant.CtxKeyAdmin).(*repository.Admins)
	if !ok {
		return errors.New("context cast error")
	}

	role, err := roleModel.FindOne(r.Context(), admin.RoleId)
	if err != nil {
		return err
	}

	if !check(role) {
		err = errors.New("authorization error")
		return errorx.New(err, errorx.Forbidden)
	}

	return nil
}
