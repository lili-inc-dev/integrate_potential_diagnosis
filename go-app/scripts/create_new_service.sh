#!/bin/bash

# input service name
read -p "Input service name(ex: example): " service_name
read -p "Input service upper name(ex: Example): " upper_service_name

# create tempalate  files
mkdir -p services/$service_name
cd services/$service_name

# create proto/common.proto
echo "syntax = \"v1\"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: // TODO: add author
)


service $service_name-api {
}" >$service_name.api

# generate template directories and
mkdir -p internal/repository internal/api internal/util database/migrations database/seeder/master database/seeder/seeder

echo -e "How to Create Basic CRUD\n\tBasic CRUD processing can be created automatically by defining a sql.\n\tTo generate repository code,Please use ./scripts/create_template.sh"
