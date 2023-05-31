package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/80andCo/LiLi-LABO/services/lili/database"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/errorx"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/handler"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/repository"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/lili-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	if err := loadConfig(&c); err != nil {
		log.Panicf("failed to load config: %#v", err)
	}

	if c.AppEnv != "local" {
		database.RunDBMigration(c.MigrationURL, "mysql://"+c.DataSource)
	}

	handler.InitCookieStore(c.SessionKey)

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors(c.FrontURL, c.FrontAdminURL),
	)
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		log.Panicf("failed to initialize ServiceContext: %#v", err)
	}

	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(errHandler)

	// serve static files
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/static/:file",
		Handler: http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP,
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func loadConfig(c *config.Config) error {
	if err := conf.Load(*configFile, c); err != nil {
		return err
	}

	if err := os.Setenv("TZ", c.TZ); err != nil {
		return err
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", c.GoogleApplicationCredentials); err != nil {
		return err
	}

	if err := os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", c.FirebaseAuthEmulatorHost); err != nil {
		return err
	}

	if err := os.Setenv("GCLOUD_PROJECT", c.GCloudProject); err != nil {
		return err
	}

	return nil
}

type responseBody struct {
	AppErrCode string `json:"app_error_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func errHandler(err error) (int, interface{}) {
	logx.Errorf("%+v", err)

	err = repository.ToAppError(err)

	body := responseBody{
		Message: err.Error(),
	}

	var appErr *errorx.AppErr
	if errors.As(err, &appErr) {
		body.AppErrCode = appErr.Code.String()
	} else {
		body.AppErrCode = errorx.UnknownError.String()
	}

	return http.StatusBadRequest, body
}
