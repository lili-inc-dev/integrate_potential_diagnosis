package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
)

var envFile = ".env"
var db *sql.DB
var finished []string
var configFile = flag.String("f", "etc/lili-api.yaml", "the config file")

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("you need to specify seeder directory")
	}
	err := loadEnv()

	if err != nil {
		panic(err)
	}

	db, err = sql.Open("mysql", dsn())
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	if err := runSeed(); err != nil {
		panic(err)
	}
}

// .env ファイルが存在する場合は読み込む
// 既に環境変数が存在する場合は上書きされないので注意
// https://github.com/joho/godotenv/blob/e74c6cadd5d7f26640f54278dc2ac083d639c505/godotenv.go#L41
func loadEnv() error {
	_, err := os.Stat(envFile)
	if err == nil {
		return godotenv.Load(envFile)
	}

	// ファイルが存在しなくてもいいので nil を返す
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// ファイルが存在しない以外のエラーは想定外
	return err
}

func runSeed() error {
	for _, dir := range os.Args[1:] {
		if err := filepath.Walk(dir, walk); err != nil {
			return err
		}

	}

	return nil
}

func walk(path string, info fs.FileInfo, err error) error {
	if alreadyFinished(path) {
		// 既に実行済みのものは、実行しない
		return nil
	}

	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	sql, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return fmt.Errorf("sql error\nsql file: %s\nerror: %w", path, err)
	}

	finished = append(finished, path)
	return nil
}

func alreadyFinished(path string) bool {
	for _, p := range finished {
		if p == path {
			return true
		}
	}
	return false
}

func dsn() string {
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// ref: https://github.com/go-sql-driver/mysql#dsn-data-source-name
	return c.DataSource
}
