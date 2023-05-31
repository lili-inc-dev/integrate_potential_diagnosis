package external

import (
	"context"
	"net/mail"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Email interface {
	Send(ctx context.Context, subject string, content string, emails []string) error
}

type email struct {
	AwsConfig           aws.Config
	ServiceEmailAddress string
}

func NewEmail(ctx context.Context, cfg config.Config) (e Email, err error) {
	defer func() {
		err = errors.Wrap(err, "NewEmail error")
	}()

	awsConfig, err := newAwsConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &email{
		AwsConfig:           awsConfig,
		ServiceEmailAddress: cfg.ServiceEmailAddress,
	}, nil
}

func newAwsConfig(ctx context.Context, cfg config.Config) (c aws.Config, err error) {
	defer func() {
		err = errors.Wrap(err, "newAwsConfig error")
	}()

	awsConf, err := AwsConfig.LoadDefaultConfig(
		ctx,
		AwsConfig.WithRegion("ap-northeast-1"),
		AwsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AwsAccessKey,
				cfg.AwsSecret,
				"",
			),
		),
	)
	return awsConf, err
}

func (e *email) Send(ctx context.Context, subject string, content string, emails []string) (err error) {
	defer func() {
		err = errors.Wrap(err, "Send error")
	}()

	ses := sesv2.NewFromConfig(e.AwsConfig)

	fromAddress := mail.Address{Name: constant.ServiceName, Address: e.ServiceEmailAddress}
	fromAddressStr := fromAddress.String()

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &fromAddressStr,
		Destination: &types.Destination{
			ToAddresses: emails,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: &content,
					},
				},
				Subject: &types.Content{
					Data: &subject,
				},
			},
		},
	}

	// メール送信
	_, err = ses.SendEmail(ctx, input)

	return err
}
