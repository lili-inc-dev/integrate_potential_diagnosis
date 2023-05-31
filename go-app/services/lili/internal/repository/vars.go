package repository

import (
	"github.com/pkg/errors"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

func ToAppError(err error) error {
	if errors.Is(err, ErrNotFound) {
		return errorx.New(err, errorx.ResourceNotFound)
	}
	return err
}
