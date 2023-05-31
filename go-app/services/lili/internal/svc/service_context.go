package svc

import (
	"context"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/external"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/middleware"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                           config.Config
	AdminModel                       repository.AdminsModel
	UserModel                        repository.UsersModel
	InactiveUserModel                repository.InactiveUsersModel
	AdminRoleModel                   repository.AdminRolesModel
	MarketValueModel                 repository.MarketValuesModel
	MarketValueDiagnoseQuestionModel repository.MarketValueDiagnoseQuestionsModel
	MarketValueDiagnoseChoiceModel   repository.MarketValueDiagnoseChoicesModel
	MarketValueDiagnoseAnswerModel   repository.MarketValueDiagnoseAnswersModel
	AuthenticateAdmin                rest.Middleware
	AuthenticateUser                 rest.Middleware
	PersonalityDiagnoseChoiceModel   repository.PersonalityDiagnoseChoicesModel
	PersonalityDiagnoseQuestionModel repository.PersonalityDiagnoseQuestionsModel
	PersonalityDiagnoseAnswerModel   repository.PersonalityDiagnoseAnswersModel
	GenderModel                      repository.GendersModel
	InterestTopicModel               repository.InterestTopicsModel
	UserInterestTopicModel           repository.UserInterestTopicsModel
	ServiceTriggerModel              repository.ServiceTriggersModel
	GeneralUserProfileModel          repository.GeneralUserProfilesModel
	UserAddressesModel               repository.UserAddressesModel
	DesiredAnnualIncomeModel         repository.DesiredAnnualIncomesModel
	LineAccountModel                 repository.LineAccountsModel
	NoticeModel                      repository.NoticesModel
	CareerWorkAnswerModel            repository.CareerWorkAnswersModel
	DiagnoseModel                    repository.DiagnosesModel
	UserCsStatusModel                repository.UserCsStatusesModel
	UnregistrationReasonModel        repository.UnregistrationReasonsModel
	UnregisteredUserProfileModel     repository.UnregisteredUserProfilesModel
	EmailAuthenticationCodeModel     repository.EmailAuthenticationCodesModel
	UserTypeModel                    repository.UserTypesModel
	FirebaseAuth                     external.FirebaseAuth
	LineAPI                          external.LineAPI
	Email                            external.Email
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	adminModel := repository.NewAdminsModel(sqlx.NewMysql(c.DataSource), c.Cache)
	userModel := repository.NewUsersModel(sqlx.NewMysql(c.DataSource), c.Cache)

	fAuth, err := external.NewFirebaseAuth(c)
	if err != nil {
		return nil, err
	}

	authAdmin := middleware.NewAuthenticateAdminMiddleware(
		adminModel,
		fAuth,
	).Handle

	authUser := middleware.NewAuthenticateUserMiddleware(
		userModel,
		fAuth,
	).Handle

	lineAPI, err := external.NewLineAPI(c)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	email, err := external.NewEmail(ctx, c)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:                           c,
		AuthenticateAdmin:                authAdmin,
		AuthenticateUser:                 authUser,
		AdminModel:                       adminModel,
		UserModel:                        repository.NewUsersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		InactiveUserModel:                repository.NewInactiveUsersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		AdminRoleModel:                   repository.NewAdminRolesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		MarketValueModel:                 repository.NewMarketValuesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		MarketValueDiagnoseQuestionModel: repository.NewMarketValueDiagnoseQuestionsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		MarketValueDiagnoseChoiceModel:   repository.NewMarketValueDiagnoseChoicesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		MarketValueDiagnoseAnswerModel:   repository.NewMarketValueDiagnoseAnswersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		PersonalityDiagnoseChoiceModel:   repository.NewPersonalityDiagnoseChoicesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		PersonalityDiagnoseQuestionModel: repository.NewPersonalityDiagnoseQuestionsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		PersonalityDiagnoseAnswerModel:   repository.NewPersonalityDiagnoseAnswersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		GenderModel:                      repository.NewGendersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		InterestTopicModel:               repository.NewInterestTopicsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UserInterestTopicModel:           repository.NewUserInterestTopicsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		ServiceTriggerModel:              repository.NewServiceTriggersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		GeneralUserProfileModel:          repository.NewGeneralUserProfilesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UserAddressesModel:               repository.NewUserAddressesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		DesiredAnnualIncomeModel:         repository.NewDesiredAnnualIncomesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		LineAccountModel:                 repository.NewLineAccountsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		NoticeModel:                      repository.NewNoticesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		CareerWorkAnswerModel:            repository.NewCareerWorkAnswersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		DiagnoseModel:                    repository.NewDiagnosesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UserCsStatusModel:                repository.NewUserCsStatusesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UnregistrationReasonModel:        repository.NewUnregistrationReasonsModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UnregisteredUserProfileModel:     repository.NewUnregisteredUserProfilesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		EmailAuthenticationCodeModel:     repository.NewEmailAuthenticationCodesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		UserTypeModel:                    repository.NewUserTypesModel(sqlx.NewMysql(c.DataSource), c.Cache),
		FirebaseAuth:                     fAuth,
		LineAPI:                          lineAPI,
		Email:                            email,
	}, nil
}
