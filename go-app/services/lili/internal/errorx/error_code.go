package errorx

type AppErrCode string

const (
	ResourceNotFound                      AppErrCode = "resource_not_found"
	InvalidParameter                      AppErrCode = "invalid_parameter"
	EmailAlreadyUsed                      AppErrCode = "email_already_used"
	ReachedEmailAuthCodeCheckAttemptLimit AppErrCode = "reached_email_auth_code_check_attempt_limit"
	WrongEmailAuthCode                    AppErrCode = "wrong_email_auth_code"
	EmailAuthCodeExpired                  AppErrCode = "email_auth_code_expired"
	Unauthorized                          AppErrCode = "unauthorized"
	Forbidden                             AppErrCode = "forbidden"
	UnknownError                          AppErrCode = "unknown_error"
)

func (c AppErrCode) String() string {
	return string(c)
}
