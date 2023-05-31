#!/bin/bash

#!/bin/bash

echo "Input values >"
declare -a lines=()
while read line; do
		if [ "$line" = "" ]; then
				break
		fi
		lines+=("$line")
done

echo "Output >"
n=${#lines[@]}
for (( i=0; i<n; i++)); do
	values=(${lines[i]})
	j=$((i+1))
	echo "($j, 1, ${values[0]}),"
	echo "($j, 2, ${values[1]}),"
done
