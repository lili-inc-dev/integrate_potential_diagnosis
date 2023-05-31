package util

import (
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"
)

// [min, max) の範囲で乱数を生成する
func GenerateRandomNumber(min, max uint64) (n uint64, err error) {
	defer func() {
		err = errors.Wrap(err, "GenerateRandomNumber error")
	}()

	if min > max {
		return 0, errors.New("`min` must not be greater than `max`")
	}

	diff := int64(max - min)
	nBig, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}

	n = uint64(nBig.Int64()) + min
	return n, nil
}
