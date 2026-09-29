#!/usr/bin/bash

# This script invokes the Docker API to create containers (acting as bots)

COUNT=${1:-5} # Takes in user input for N bots
for ((i=1; i<=COUNT; i++)); do
    NAME=$(printf "bot_%03d" "$i")
    echo "[+] Starting $NAME"
    
    # Remove any existing container
    if sudo docker container inspect "$NAME" &>/dev/null; then
        sudo docker rm -f "$NAME" >/dev/null
    fi

    sudo docker run -d \
    --name "$NAME" \
    --add-host=host.docker.internal:host-gateway \
    bot_image \
    --master ws://host.docker.internal:8080/connect

done