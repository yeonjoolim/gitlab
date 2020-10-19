#!/bin/bash

a="\"http://localhost/api/v4/projects/1/registry/repositories/1/tags/${1:27}\""
b="curl --request DELETE --header \"PRIVATE-TOKEN: Kca5ANW3gda32Jn3M6C9\" $a"

cd ./clair
clairctl-linux-amd64 --config=clairctl.yml analyze $1>result.txt
clairctl-linux-amd64 --config=clairctl.yml report $1
python filter.py
python parse.py>score.txt
file="score.txt"

while IFS= read -r line
do
	if [ "$line" == "Delete your docker image" ]; then
		echo "-----<Alert> This Image is alot Vulnerability detection-----"
		$c='docker rmi '$1
		echo "Remove Image on GitLab"
		echo $b > Del.sh
        	chmod +x Del.sh
        	result=`./Del.sh`
        	echo " " $result
	fi
done < "$file"

cat score.txt
