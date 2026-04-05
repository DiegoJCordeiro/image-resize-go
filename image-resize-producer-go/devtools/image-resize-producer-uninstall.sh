#!/usr/bin/env sh

aws s3 rm --recursive s3://diegojc-bucket/

docker compose down -v --rmi all
