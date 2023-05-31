package util

import (
	"time"

	"github.com/pkg/errors"
)

func StringToTime(strintTime string) (t time.Time, err error) {
	defer func() {
		err = errors.Wrap(err, "StringToTime error")
	}()

	return time.Parse(time.RFC3339, strintTime)
}

func DateStrToTime(strintTime string) (t time.Time, err error) {
	defer func() {
		err = errors.Wrap(err, "DateStrToTime error")
	}()

	return time.Parse("2006/01/02", strintTime)
}

func TimeToStringFormatDate(toStrTime time.Time) string {
	return toStrTime.Format("2006/01/02")
}

func TimeToStringFormatDateJp(toStrTime time.Time) string {
	return toStrTime.Format("2006年1月2日")
}

func TimeToStringFormatJp(toStrTime time.Time) string {
	return toStrTime.Format("2006年1月2日 15:04")
}
