package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
)

const (
	claimKeyIsValid = "is_valid"
)

func newFirebaseApp() (app *firebase.App, err error) {
	defer func() {
		err = errors.Wrap(err, "newFirebaseApp error")
	}()

	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	return app, nil
}

func newFirebaseAuthClient(cfg config.Config) (client *auth.Client, err error) {
	defer func() {
		err = errors.Wrap(err, "newFirebaseAuthClient error")
	}()

	app, err := newFirebaseApp()
	if err != nil {
		return nil, err
	}

	if cfg.AppEnv == "test" && cfg.FirebaseAuthEmulatorHost == "" {
		return nil, errors.New("use Firebase Auth Emulator in test environment")
	}

	client, err = app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase auth client: %w", err)
	}

	return client, nil
}

type (
	FirebaseAuth interface {
		FindUser(ctx context.Context, firebaseUID string) (*auth.UserRecord, error)
		FindUserList(ctx context.Context, firebaseUIDList []string) ([]*auth.UserRecord, error)
		FindUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
		VerifyIDToken(ctx context.Context, idToken string) (token *auth.Token, err error)
		SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) error
		CreateUser(ctx context.Context, name string) (*auth.UserRecord, error)
		CreateUserWithEmailAndPassword(ctx context.Context, name string, email string, password string) (*auth.UserRecord, error)
		CreateCustomToken(ctx context.Context, fUID string) (token string, err error)
		ChangePassword(ctx context.Context, uid string, password string) error
		ChangeEmail(ctx context.Context, uid, email string) error
		UpdateEmailVerified(ctx context.Context, uid string, emailVerified bool) (err error)
		DeleteUser(ctx context.Context, uid string) error
		ConfirmEmail(ctx context.Context, actionCode string) (email string, err error)
		GenerateEmailVerificationLink(ctx context.Context, email string) (string, error)
		RevokeRefreshTokens(ctx context.Context, uid string) error
	}

	firebaseAuth struct {
		cl  *auth.Client
		cfg config.Config
	}
)

func NewFirebaseAuth(cfg config.Config) (auth FirebaseAuth, err error) {
	defer func() {
		err = errors.Wrap(err, "NewFirebaseAuth error")
	}()

	cl, err := newFirebaseAuthClient(cfg)
	if err != nil {
		return nil, err
	}

	return &firebaseAuth{
		cl:  cl,
		cfg: cfg,
	}, nil
}

// firebase auth apiに必要なクエリ文字列をセット
func (f *firebaseAuth) setFirebaseAuthAPIQueryParam(params *url.Values) {
	params.Add("key", f.cfg.FirebaseWebAPIKey)
}

func (f *firebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (token *auth.Token, err error) {
	defer func() {
		err = errors.Wrap(err, "VerifyIDToken error")
	}()

	token, err = f.cl.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	claims := token.Claims
	value, ok := claims[claimKeyIsValid]
	if !ok {
		return nil, errors.New("firebase token claim is invalid")
	}
	if !value.(bool) {
		return nil, errors.New("invalid firebase user")
	}

	return token, nil
}

// 既存のカスタムクレイムをclaimsで置き換える
// 既存クレイムの内claimsに含まれないものは消える
func (f *firebaseAuth) SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) (err error) {
	defer func() {
		err = errors.Wrap(err, "SetCustomClaims error")
	}()

	claims[claimKeyIsValid] = true
	return f.cl.SetCustomUserClaims(ctx, uid, claims)
}

func (f *firebaseAuth) CreateCustomToken(ctx context.Context, fUID string) (token string, err error) {
	defer func() {
		err = errors.Wrap(err, "CreateCustomToken error")
	}()

	return f.cl.CustomToken(ctx, fUID)
}

func (f *firebaseAuth) FindUser(ctx context.Context, firebaseUID string) (user *auth.UserRecord, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUser error")
	}()

	user, err = f.cl.GetUser(ctx, firebaseUID)
	if auth.IsUserNotFound(err) {
		return nil, errorx.New(err, errorx.ResourceNotFound)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (f *firebaseAuth) FindUserList(ctx context.Context, firebaseUIDList []string) (users []*auth.UserRecord, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserList error")
	}()

	identifiers := []auth.UserIdentifier{}
	for _, uid := range firebaseUIDList {
		identifiers = append(identifiers, auth.UIDIdentifier{UID: uid})
	}

	result, err := f.cl.GetUsers(ctx, identifiers)
	if err != nil {
		return nil, err
	}
	if len(result.NotFound) > 0 {
		err := fmt.Errorf("firebaseUIDListに存在しないUIDが含まれます: %#v", result.NotFound)
		return nil, errorx.New(err, errorx.ResourceNotFound)
	}

	return result.Users, nil
}

func (f *firebaseAuth) FindUserByEmail(ctx context.Context, email string) (user *auth.UserRecord, err error) {
	defer func() {
		err = errors.Wrap(err, "FindUserByEmail error")
	}()

	user, err = f.cl.GetUserByEmail(ctx, email)
	if auth.IsUserNotFound(err) {
		return nil, errorx.New(err, errorx.ResourceNotFound)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (f *firebaseAuth) CreateUser(
	ctx context.Context,
	name string,
) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		DisplayName(name).
		Disabled(false)

	user, err := f.cl.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	claims := map[string]interface{}{claimKeyIsValid: true}
	if err := f.cl.SetCustomUserClaims(ctx, user.UID, claims); err != nil {
		return nil, err
	}
	user.CustomClaims = claims

	return user, nil
}

func (f *firebaseAuth) CreateUserWithEmailAndPassword(
	ctx context.Context,
	name string,
	email string,
	password string,
) (user *auth.UserRecord, err error) {
	defer func() {
		err = errors.Wrap(err, "CreateUser error")
	}()

	params := (&auth.UserToCreate{}).
		DisplayName(name).
		Email(email).
		EmailVerified(false).
		Password(password).
		Disabled(false)

	user, err = f.cl.CreateUser(ctx, params)
	if auth.IsEmailAlreadyExists(err) {
		return nil, errorx.New(err, errorx.EmailAlreadyUsed)
	}
	if err != nil {
		return nil, err
	}

	claims := map[string]interface{}{claimKeyIsValid: true}
	if err := f.cl.SetCustomUserClaims(ctx, user.UID, claims); err != nil {
		return nil, err
	}
	user.CustomClaims = claims

	return user, nil
}

func (f *firebaseAuth) ChangePassword(
	ctx context.Context,
	uid string,
	password string,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "ChangePassword error")
	}()

	params := (&auth.UserToUpdate{}).
		Password(password)

	_, err = f.cl.UpdateUser(ctx, uid, params)

	return err
}

func (f *firebaseAuth) ChangeEmail(ctx context.Context, uid, email string) (err error) {
	defer func() {
		err = errors.Wrap(err, "ChangeEmail error")
	}()

	params := (&auth.UserToUpdate{}).
		Email(email).
		EmailVerified(false)

	_, err = f.cl.UpdateUser(ctx, uid, params)
	if auth.IsEmailAlreadyExists(err) {
		return errorx.New(err, errorx.EmailAlreadyUsed)
	}

	return err
}

func (f *firebaseAuth) UpdateEmailVerified(ctx context.Context, uid string, emailVerified bool) (err error) {
	defer func() {
		err = errors.Wrap(err, "UpdateEmailVerified error")
	}()

	params := (&auth.UserToUpdate{}).
		EmailVerified(emailVerified)

	_, err = f.cl.UpdateUser(ctx, uid, params)
	return err
}

func (f *firebaseAuth) DeleteUser(ctx context.Context, uid string) (err error) {
	defer func() {
		err = errors.Wrap(err, "DeleteUser error")
	}()

	return f.cl.DeleteUser(ctx, uid)
}

// https://firebase.google.com/docs/reference/rest/auth#section-confirm-email-verification
func (f *firebaseAuth) ConfirmEmail(ctx context.Context, actionCode string) (email string, err error) {
	defer func() {
		err = errors.Wrap(err, "ConfirmEmail error")
	}()

	reqBody := map[string]string{
		"oobCode": actionCode,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", errorx.New(err, errorx.InvalidParameter)
	}

	req, err := http.NewRequest(
		"POST",
		config.FirebaseConfirmEmailVerificationURL(f.cfg.FirebaseAuthEmulatorHost),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	params := req.URL.Query()
	f.setFirebaseAuthAPIQueryParam(&params)
	req.URL.RawQuery = params.Encode()

	util.SetReqContentTypeJSON(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("faild to confirm email\nstatus code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resp struct {
		Email string `json:"email"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.Email, nil
}

func (f *firebaseAuth) GenerateEmailVerificationLink(ctx context.Context, email string) (link string, err error) {
	defer func() {
		err = errors.Wrap(err, "GenerateEmailVerificationLink error")
	}()

	link, err = f.cl.EmailVerificationLinkWithSettings(ctx, email, nil)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (f *firebaseAuth) RevokeRefreshTokens(ctx context.Context, uid string) (err error) {
	defer func() {
		err = errors.Wrap(err, "RevokeRefreshTokens error")
	}()

	if err := f.cl.RevokeRefreshTokens(ctx, uid); err != nil {
		return err
	}
	return nil
}
