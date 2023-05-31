package constant

const (
	FirebaseConfirmEmailVerificationEndpoint = "identitytoolkit.googleapis.com/v1/accounts:update"

	FirebaseCustomClaimKeyLineID = "line_id"

	FirebaseCustomClaimKeySignUpState                = "sign_up_state"  // 登録状態を表すカスタムクレイムのキー
	FirebaseCustomClaimValueSignUpStateEmailVerified = "email_verified" // メールアドレス確認済
	FirebaseCustomClaimValueSignUpStateRegisterd     = "registered"     // 本登録済
)
