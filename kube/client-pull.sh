#!/bin/bash

docker login gitlab.pel

./tls-pull-client $1
if [ $? -eq 0 ];then
    docker pull $1
    echo "-----------Success Pull Image-----------"
    if [ $? -eq 0 ];then
        echo "-----------Success receive Image Sign data-----------"
    	./layer-verify $1
	if [ $? -eq 0 ];then
        	echo "-----------Success Verify Image-----------"
	else
		echo "-----------Fail Verify, Remove Image----------"
		b=$(docker ps -a | grep $1 | awk '{print $1}')
		if [ ! -z "$b" ];then
        		docker rm $b
		fi
		docker rmi $1
	fi
    else
	echo "Receive Fail Sign data"
	c=$(docker ps -a | grep $1 | awk '{print $1}')
	if [ ! -z "$c" ];then
        	docker rm $c
	fi
	docker rmi $1	
    fi
else
   echo "Fail Pull Image"
fi

