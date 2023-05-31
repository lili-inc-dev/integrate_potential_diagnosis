package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunDBMigration(migrationURL string, dbSource string) {
	fmt.Println(migrationURL, dbSource)
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("migration err: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}
	log.Println("migration done")
}
