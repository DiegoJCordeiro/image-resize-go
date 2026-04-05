#!/usr/bin/env sh

docker-compose down
docker-compose up -d --force-recreate

aws s3 mb s3://diegojc-bucket