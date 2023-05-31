package logic

import (
	"context"
	"fmt"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/types"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ContactLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContactLogic {
	return &ContactLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContactLogic) Contact(req *types.ContactReq) (err error) {
	defer func() {
		err = errors.Wrap(err, "Contact error")
	}()

	var reqContent string
	if req.Content != nil {
		reqContent = *req.Content
	} else {
		reqContent = ""
	}

	content := fmt.Sprintf(`
お客様名：
%s

お問い合わせ内容:
%s

メールアドレス:
%s
`, req.Name, reqContent, req.Email)

	if err := l.svcCtx.Email.Send(l.ctx, req.Name+"さんからお問い合わせがありました", content, []string{l.svcCtx.Config.ContactEmail}); err != nil {
		return err
	}
	return nil
}
