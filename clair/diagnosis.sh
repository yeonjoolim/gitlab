#!/bin/bash

cd ./clair
clairctl-linux-amd64 --config=clairctl.yml analyze $1>result.txt
clairctl-linux-amd64 --config=clairctl.yml report $1
python filter.py
python parse.py>score.txt
file="score.txt"
while IFS= read -r line
do
	if [ "$line" == "Delete your a docker image" ]; then
		docker rmi $1 
	fi
done < "$file"

cat score.txt
