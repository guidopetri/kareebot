#!/bin/bash
git pull
touch doods.txt
docker build -t kareebot .
container=`docker ps -aqf "name=kareebot"`
$(docker stop $container)
$(docker rm $container)
docker run --name kareebot -d kareebot
