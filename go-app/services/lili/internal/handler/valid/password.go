package valid

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	passwordMinLength = 8
)

var (
	containsLowerRegex  = regexp.MustCompile(`^.*[a-z].*$`)
	containsUpperRegex  = regexp.MustCompile(`^.*[A-Z].*$`)
	containsNumberRegex = regexp.MustCompile(`^.*[0-9].*$`)
	containsSymbolRegex = regexp.MustCompile(`^.*[!"#$%&'()*+,-./:;<=>?@[\]^_` + "`" + `{|}~].*$`)

	/*
	  半角英数字記号のみ
	*/
	passwordRegex = regexp.MustCompile(`^[!-~]+$`)
)

/*
パスワード要件
- [8]文字以上
- 英小文字、英大文字、数字、記号のうち3種類以上の文字種の組み合わせ
*/
func ValidatePassword(pw string) error {
	cnt := utf8.RuneCountInString(pw)
	if cnt < passwordMinLength {
		return fmt.Errorf("パスワードは%d文字以上にしてください", passwordMinLength)
	}

	if !passwordRegex.MatchString(pw) {
		return errors.New("password validation error: 英小文字、英大文字、数字、記号以外が含まれています")
	}

	regexes := []*regexp.Regexp{containsLowerRegex, containsUpperRegex, containsNumberRegex, containsSymbolRegex}
	matchedCnt := 0
	for _, regex := range regexes {
		if regex.MatchString(pw) {
			matchedCnt++
		}
	}

	if matchedCnt < 3 {
		return errors.New("password validation error: 英小文字、英大文字、数字、記号のうち3種類以上の文字種を組み合わせてください")
	}

	return nil
}
