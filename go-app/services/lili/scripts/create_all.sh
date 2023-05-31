#!/bin/bash

## Get variables
service_name=lili
upper_service_name=Lili
read -p "Input lower model name: " lower_model_name
read -p "Input upper model name: " upper_model_name

# Select Wheter create new sql file or not
read -p "Do you want to create new migration files (y/n) " yn
case $yn in
# Create sql files
y) migrate create -ext sql -dir database/migrations $lower_model_name ;;
n) echo "We skipped creating migration files" ;;
*)
	echo invalid response
	exit 1
	;;

# Create repository files from _.up.sql
echo "Please fill generated sql files."
read -p "target sql file name to generate repository: " sql_file_name
goctl model mysql ddl -c -style go_zero -src database/migrations/$sql_file_name -dir internal/repository

# Get example params for lili.api
go run ../../pkg/go-protogen/main.go api --inputFilePath internal/repository/"$lower_model_name"s_model_gen.go --model $upper_model_name

echo "You can use auto-generated code that is displayed above."
read -p "Did you renew $service_name.api file? (y/n) " yn
case $yn in
y) echo ok, we will generate template code ;;
n)
	echo stopping...
	exit
	;;
*)
	echo invalid response
	exit 1
	;;
esac

# Make changes to the configuration files
gsed -i "/service $service_name-api/a \ \t@server(\n\t\thandler: Create"$upper_model_name"Handler\n\t)\n\tpost /"$lower_model_name"(Create"$upper_model_name"Req)\n\t@server(\n\t\thandler: Update"$upper_model_name"Handler\n\t)\n\tput /"$lower_model_name"(Update"$upper_model_name"Req)\n\t@server(\n\t\thandler: Delete"$upper_model_name"Handler\n\t)\n\tdelete /"$lower_model_name"(Delete"$upper_model_name"Req)\n\t@server(\n\t\thandler: FindById"$upper_model_name"Handler\n\t)\n\tget /"$lower_model_name"(FindById"$upper_model_name"Req) returns(FindById"$upper_model_name"Res)" $service_name.api
gsed -i "/DataSource string/a \ \t"$upper_model_name"Table string" internal/config/config.go
gsed -i "/Config     config.Config/a \ \t"$upper_model_name"Model repository."$upper_model_name"sModel" internal/svc/service_context.go
gsed -i "/Config:     c/a \ \t"$upper_model_name"Model: repository.New"$upper_model_name"sModel(sqlx.NewMysql(c.DataSource), c.Cache)," internal/svc/service_context.go
gsed -i "/DataSource/a"$upper_model_name"Table: "$lower_model_name"s" etc/lili-api.yaml


# Create template files(handler,logic)
goctl api go -api $service_name.api -style go_zero -dir .

# generate swagger file
goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api lili.api -dir ./docs/

