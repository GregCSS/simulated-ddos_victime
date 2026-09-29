#!/usr/bin/bash

# This script deletes all containers created by create_botnet.sh

docker ps -a --filter "name=bot_" --format "{{.Names}}" | while read -r container; do
    docker rm -f "$container"
done