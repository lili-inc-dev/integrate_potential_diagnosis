package errorx

import (
	"fmt"

	"github.com/pkg/errors"
)

type Err error

type AppErr struct {
	Err
	Code AppErrCode
}

func New(err error, code AppErrCode) *AppErr {
	return &AppErr{Err: err, Code: code}
}

func HasAppErrCode(err error, code AppErrCode) bool {
	var appErr *AppErr
	if !errors.As(err, &appErr) {
		return false
	}

	return appErr.Code == code
}

func (a *AppErr) Format(s fmt.State, verb rune) {
	if f, ok := a.Err.(fmt.Formatter); ok {
		f.Format(s, verb)
		return
	}

	fmt.Fprintf(s, "%+v", a.Err)
}
