package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
)

func GenerateUlid() (u ulid.ULID, err error) {
	defer func() {
		err = errors.Wrap(err, "GenerateUlid error")
	}()

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	return ulid.New(ms, entropy)
}
