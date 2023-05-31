package main

import (
	"database/sql"
	"flag"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "./etc/lili-api.yaml", "the config file")

func main() {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	db, _ := sql.Open("mysql", c.DataSource+"&multiStatements=true")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"../../database/migrations",
		"mysql",
		driver,
	)

	if err := m.Steps(2); err != nil {
		panic(err)
	}
}
