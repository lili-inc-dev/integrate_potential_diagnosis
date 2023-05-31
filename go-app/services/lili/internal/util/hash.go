package util

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(s string) (hash string, err error) {
	defer func() {
		err = errors.Wrap(err, "GenerateHash error")
	}()

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func GenerateHashOfNumber(n uint64) (hash string, err error) {
	defer func() {
		err = errors.Wrap(err, "GenerateHashOfNumber error")
	}()

	bytes := Uint64ToBytes(n)
	hashBytes, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}
