package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/external"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInactiveUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInactiveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInactiveUserLogic {
	return &CreateInactiveUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	emailAuthCodeExpirationHour = 1 // メール認証コードの有効期間（単位：時間）
)

func (l *CreateInactiveUserLogic) CreateInactiveUser(req *types.CreateInactiveUserReq, firebaseUID, lineID string) (resp *types.CreateInactiveUserRes, err error) {
	defer func() {
		err = errors.Wrap(err, "CreateInactiveUser error")
	}()

	fbUser, err := l.svcCtx.FirebaseAuth.FindUserByEmail(l.ctx, req.Email)
	if err != nil && !errorx.HasAppErrCode(err, errorx.ResourceNotFound) /* 想定外のエラー */ {
		return nil, err
	}

	if err == nil /* FindUserByEmailでfirebase userが見つかった */ {
		if signUpState, ok := fbUser.CustomClaims[constant.FirebaseCustomClaimKeySignUpState]; ok {
			ngs := []string{
				constant.FirebaseCustomClaimValueSignUpStateEmailVerified,
				constant.FirebaseCustomClaimValueSignUpStateRegisterd,
			}
			if util.Contains(ngs, signUpState.(string)) {
				err = errors.New("the provided email has been already used")
				return nil, errorx.New(err, errorx.EmailAlreadyUsed)
			}
		}

		if fbUser.UID != firebaseUID /* firebase userが自分ではない */ {
			// 見つかったfirebase userは、仮登録は終えたがメール確認は終えていないままなので削除する
			if err := l.svcCtx.FirebaseAuth.DeleteUser(l.ctx, fbUser.UID); err != nil {
				return nil, err
			}
		}
	}

	// firebase userにemailを紐付ける
	// LINEログイン時は紐付けられていない
	if err := l.svcCtx.FirebaseAuth.ChangeEmail(l.ctx, firebaseUID, req.Email); err != nil {
		return nil, err
	}

	ulid, err := util.GenerateUlid()
	if err != nil {
		return nil, err
	}

	// 以下のケースでline_idが被るのでInsertではなくUpsert
	// - 本登録を完了せず再度仮登録を行う
	// - 退会後に仮登録を行う
	inactiveUser, err := l.svcCtx.InactiveUserModel.Upsert(l.ctx, &repository.InactiveUsers{
		Id:          ulid.String(),
		LineId:      lineID,
		TypeId:      req.TypeId,
		FirebaseUid: firebaseUID,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}

	authCode, err := generateEmailAuthCode()
	if err != nil {
		return nil, err
	}

	authCodeHash, err := util.GenerateHash(authCode)
	if err != nil {
		return nil, err
	}

	oldEmailAuthCode, err := l.svcCtx.EmailAuthenticationCodeModel.FindByInactiveUserIdNoCache(l.ctx, inactiveUser.Id)
	if err != nil && !errors.Is(err, repository.ErrNotFound) /* 想定外のエラー */ {
		return nil, err
	}

	needToCreateNewEmailAuthCode := true // 新規にメール認証コードレコードを作成必要かどうか

	if err == nil /* メール認証コードレコードが存在する */ {
		expiration := oldEmailAuthCode.CreatedAt.Add(emailAuthCodeExpirationHour * time.Hour)
		isExpired := time.Now().After(expiration)
		if !isExpired {
			needToCreateNewEmailAuthCode = false
		}
	}

	if needToCreateNewEmailAuthCode {
		ulid, err := util.GenerateUlid()
		if err != nil {
			return nil, err
		}

		if _, err = l.svcCtx.EmailAuthenticationCodeModel.InsertNoCache(l.ctx, &repository.EmailAuthenticationCodes{
			Id:             ulid.String(),
			InactiveUserId: inactiveUser.Id,
			CodeHash:       authCodeHash,
			AttemptCount:   0,
		}); err != nil {
			return nil, err
		}
	} else {
		// 既存の認証コードレコードの`code_hash`のみ更新する
		if _, err = l.svcCtx.EmailAuthenticationCodeModel.UpdateCodeHashNoCache(l.ctx, oldEmailAuthCode.Id, authCodeHash); err != nil {
			return nil, err
		}
	}

	if err := sendConfirmationMail(l.ctx, l.svcCtx.Email, req.Email, authCode); err != nil {
		return nil, err
	}

	token, err := l.svcCtx.FirebaseAuth.CreateCustomToken(l.ctx, firebaseUID)
	if err != nil {
		return nil, err
	}

	// メールアドレス変更によりフロントが持つfirebase id tokenが失効するので、再ログインに使用するカスタムトークンを返却する
	return &types.CreateInactiveUserRes{
		FirebaseCustomToken: token,
	}, nil
}

// 6桁認証コード生成
// 範囲は`000000`〜`999999`
func generateEmailAuthCode() (authCodeStr string, err error) {
	defer func() {
		err = errors.Wrap(err, "generateEmailAuthCode error")
	}()

	authCode, err := util.GenerateRandomNumber(0, 1e6)
	if err != nil {
		return "", err
	}

	authCodeStr = fmt.Sprintf("%06d", authCode)
	return authCodeStr, nil
}

func sendConfirmationMail(ctx context.Context, emailExternal external.Email, email, authCode string) (err error) {
	defer func() {
		err = errors.Wrap(err, "sendConfirmationMail error")
	}()

	content := fmt.Sprintf(`
%sをご利用いただきありがとうございます。
下記の認証コードを会員登録画面に入力し、メールアドレスの確認を完了してください。

%s

認証コードの有効期限は最大%d時間です。
`, constant.ServiceName, authCode, emailAuthCodeExpirationHour)

	if err := emailExternal.Send(
		ctx,
		"メールアドレスの確認",
		content,
		[]string{email},
	); err != nil {
		return err
	}

	return nil
}
