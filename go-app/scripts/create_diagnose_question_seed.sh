#!/bin/bash

echo "Input questions >"
declare -a questions=()
while read line; do
		if [ "$line" = "" ]; then
				break
		fi
		questions+=("$line")
done

echo "Input start number >"
read start

echo "Output >"
n=${#questions[@]}
for (( i=0; i<n; i++)); do
	position=$((i+start))
  echo "($position, '${questions[i]}'),"
done
