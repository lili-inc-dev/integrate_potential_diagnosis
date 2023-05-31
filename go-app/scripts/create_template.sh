#!/bin/bash
service_name=lili
upper_service_name=Lili
read -p "Input lower model name: " lower_model_name
read -p "Input upper model name: " upper_model_name

read -p "Do you want to create new migration files (y/n) " yn

case $yn in
y) migrate create -ext sql -dir services/$service_name/database/migrations $lower_model_name ;;
n) echo "We skipped creating migration files" ;;
*)
	echo invalid response
	exit 1
	;;
esac

# create migration files

# Please input sql
echo "Please fill generated sql files."
read -p "Did you filled sql files collectly? (y/n) " yn

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

# ceate repository files
goctl model mysql ddl -c -style go_zero -src services/$service_name/database/migrations/*.up.sql -dir services/$service_name/internal/repository

go run cmd/go-protogen/main.go api --inputFilePath services/$service_name/internal/repository/"$lower_model_name"s_model_gen.go --model $upper_model_name

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

cd services/$service_name

gsed -i "/service $service_name-api/a \ \t@server(\n\t\thandler: Create"$upper_model_name"Handler\n\t)\n\tpost /"$lower_model_name"(Create"$upper_model_name"Req)\n\t@server(\n\t\thandler: Update"$upper_model_name"Handler\n\t)\n\tput /"$lower_model_name"(Update"$upper_model_name"Req)\n\t@server(\n\t\thandler: Delete"$upper_model_name"Handler\n\t)\n\tdelete /"$lower_model_name"(Delete"$upper_model_name"Req)\n\t@server(\n\t\thandler: FindById"$upper_model_name"Handler\n\t)\n\tget /"$lower_model_name"(FindById"$upper_model_name"Req) returns(FindById"$upper_model_name"Res)" $service_name.api

goctl api go -api $service_name.api -style go_zero -dir .
