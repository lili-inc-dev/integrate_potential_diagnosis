#!/bin/bash
# #!/usr/local/bin/bash

# マイグレーションファイルのバージョンが被っていないかチェックする
# 被っているとgolang-migrate実行時にエラーが起きる

cd "./services/lili/database/migrations"

declare -A map

file_names="./*"
for name in $file_names; do
	# nameの例: "000002_create_admin_role_table.up.sql"

	splited=(${name//_/ }) # split by `_`
	splited_len=${#splited[@]}

	if [ $splited_len -lt 2 ]; then
		echo "invalid file name: $name"
		exit 1
	fi

	tail=${splited[splited_len - 1]} # tailの例: "table.up.sql"

	tails=(${tail//./ })
	tails_len=${#tails[@]}

	if [ $tails_len -lt 3 ]; then
		echo "invalid file name tail: $tail"
		exit 1
	fi

	head=${splited[0]}
	key=$head${tails[1]}

	if [ "${map[$key]}" ]; then
		echo "duplicate key: $key"
		exit 1
	fi

	map[$key]=true
done

echo "Passed"