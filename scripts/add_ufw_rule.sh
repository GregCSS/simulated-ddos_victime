#!/usr/bin/bash
set -euo pipefail

# This script adds a new rule for firewall to allow local-to-docker-network bridge. It was initially dropping packets
# coming from the bridge

SUBNET="172.17.0.0/24"
PORT="8080"

echo "[+] Allowing Docker bridge → C2 port ${PORT}"
sudo ufw allow in on docker0 from "$SUBNET" to any port "$PORT" proto tcp

echo
sudo ufw status numbered