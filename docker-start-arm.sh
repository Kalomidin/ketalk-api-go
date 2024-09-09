#!/bin/sh
set -ex

echo "downloading docker images"
sudo  docker-compose-arm --env-file docker.env.local pull
echo "starting docker containers"
sudo docker-compose-arm  --env-file docker.env.local up -d