#!/bin/bash

a="\"http://localhost:12001/api/v4/projects/1/registry/repositories/1/tags/${1:27}\""
b="curl --request DELETE --header \"PRIVATE-TOKEN: tf7BMdtezrRQ-8kn4XaH\" $a"


cd ../../../clair-and-docker-notary-example/
docker-compose run --rm clair-scanner $1>../go/gitlab/v3/docker_imagetest.txt
cd ../go/gitlab/v3/
sed '/Unapproved/!d' docker_imagetest.txt>result.txt
python confirm.py>1.txt
python tx.py>$1_score.txt
sed -n 2p $1_score.txt>score.txt
file="$1_score.txt"
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
