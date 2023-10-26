#!/bin/sh
set -ex

echo "downloading docker images"
docker-compose --env-file docker.env.local pull
echo "starting docker containers"
docker-compose  --env-file docker.env.local up -d