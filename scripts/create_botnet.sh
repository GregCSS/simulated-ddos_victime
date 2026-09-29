#!/usr/bin/bash

# This script invokes the Docker API, creates a subnet and creates containers (acting as bots) within it

COUNT=${1:-5} # Takes in user input for N bots
for ((i=1; i<=COUNT; i++)); do
    ID=$(printf "bot_%03d" "$i")
    echo "[+] Starting $ID"
    
    docker run -d --name "$ID" --network botnet bot --id "$ID" \
    --add-host=host.docker.internal:host-gateway
done