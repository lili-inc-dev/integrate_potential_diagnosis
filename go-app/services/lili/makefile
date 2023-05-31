gen-api-docs: ## Generate swagger.json from lili.api
	rm -f ./docs/swagger.json
	goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api lili.api -dir ./docs/

gen-api: ## Generate API files from lili.api
	goctl api go -api lili.api -style go_zero --dir .

gen-api-with-struct: ## Generate API files with sample struct
	./scripts/create_api_template.sh

target=""
gen-repository: ## Generate repository files from database/migrations/*.up.sql.Make sure specify ex) make gen-repository target=database/migrations/000001_admin.up.sql
	echo ${target}
	goctl model mysql ddl -c -style go_zero -src ${target} -dir internal/repository


gen-all: ## Generate repository + API + swagger files from specified sql
	./scripts/create_all.sh


migrate: ## Migrate up
	docker-compose exec lili-back migrate create -ext sql -dir database/migrations ${file}

migrate-up: ## Migrate up(you can specify num.default is all sql files)
	docker-compose exec lili-back migrate -database "mysql://root:root@tcp(lili-mysql:3306)/lili_labo" -path database/migrations up ${num}

migrate-down: ## Migration down(you can specify num.default is all sql files)
	docker-compose exec lili-back sh -c 'yes | migrate -database "mysql://root:root@tcp(lili-mysql:3306)/lili_labo" -path database/migrations down ${num}'

migrate-down-force: ## Migration force down(all)
	docker-compose exec lili-mysql sh -c 'echo "DROP DATABASE IF EXISTS lili_labo; CREATE DATABASE lili_labo;" | mysql -u root -proot'

migrate-reset: migrate-down-force migrate-up seed-all ## Reset migration and seed

seed-all: ## Seed all
	docker-compose exec lili-back go run database/seeds/seeder.go database/seeds/master database/seeds/transaction

seed-master: ## Seed master
	docker-compose exec lili-back go run database/seeds/seeder.go database/seeds/master



help: ## Display help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


.PHONY: gen-api-docs gen-api gen-api-with-struct gen-repository gen-all migrate migrate.up migrate.down migrate.down.force migrate.reset seed.all seed.master help
