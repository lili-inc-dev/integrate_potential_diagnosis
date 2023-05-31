package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DataSource                   string
	MigrationURL                 string
	AdminTable                   string
	UserTable                    string
	AdminRoleTable               string
	Cache                        cache.CacheConf
	AppEnv                       string
	FrontURL                     string
	FrontAdminURL                string
	SessionKey                   string
	TZ                           string
	GoogleApplicationCredentials string
	FirebaseAuthEmulatorHost     string
	GCloudProject                string
	FirebaseWebAPIKey            string
	LineLoginChannelID           uint64
	LineLoginChannelSecret       string
	LineMsgAPIChannelSecret      string
	LineMsgAPIChannelAccessToken string
	AwsAccessKey                 string
	AwsSecret                    string
	ServiceEmailAddress          string
	ContactEmail                 string
}
