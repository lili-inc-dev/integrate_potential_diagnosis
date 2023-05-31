#!/bin/bash

cd "../services/lili/database/migrations"

files="./*"
i=0
for f in $files; do
	names=(${f//_/ })
	version="$(((i+2) / 2))"
	new_name=$(printf "%06d\n" "${version}")

	names_len=${#names[@]}
	for (( j=1; j<names_len; j++)); do
		new_name+="_${names[j]}"
	done

	cp $f $new_name

	i=$((i+1))
done
