#!/bin/bash
service_name=lili
upper_service_name=Lili
read -p "Input lower model name: " lower_model_name
read -p "Input upper model name: " upper_model_name

go run ../../cmd/go-protogen/main.go api --inputFilePath internal/repository/"$lower_model_name"s_model_gen.go --model $upper_model_name

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

gsed -i "/service $service_name-api/a \ \t@server(\n\t\thandler: Create"$upper_model_name"Handler\n\t)\n\tpost /"$lower_model_name"(Create"$upper_model_name"Req)\n\t@server(\n\t\thandler: Update"$upper_model_name"Handler\n\t)\n\tput /"$lower_model_name"(Update"$upper_model_name"Req)\n\t@server(\n\t\thandler: Delete"$upper_model_name"Handler\n\t)\n\tdelete /"$lower_model_name"(Delete"$upper_model_name"Req)\n\t@server(\n\t\thandler: FindById"$upper_model_name"Handler\n\t)\n\tget /"$lower_model_name"(FindById"$upper_model_name"Req) returns(FindById"$upper_model_name"Res)" $service_name.api

goctl api go -api $service_name.api -style go_zero -dir .

goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api lili.api -dir ./docs/
